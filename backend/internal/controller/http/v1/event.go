package v1

import (
	"encoding/json"
	"net/http"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EventHandler struct {
	uc usecase.EventUseCase
}

func NewEventHandler(uc usecase.EventUseCase) *EventHandler {
	return &EventHandler{uc: uc}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var req request.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	wantParticipate := true
	if req.WantParticipate != nil {
		wantParticipate = *req.WantParticipate
	}

	input := dto.CreateEventInput{
		Title:           req.Title,
		Description:     req.Description,
		Rules:           req.Rules,
		Recommendations: req.Recommendations,
		OrganizerNotes:  req.OrganizerNotes,
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		EndDate:         req.EndDate,
		MaxParticipants: req.MaxParticipants,
		Budget:          req.Budget,
		WantParticipate: wantParticipate,
	}

	event, err := h.uc.Create(r.Context(), input, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.EventToResponse(&event))
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	event, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.EventToResponse(event))
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	events, err := h.uc.GetMyEvents(r.Context(), userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.EventsToResponse(events))
}

// UpdateEvent — только организатор может редактировать (проверка в usecase)
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	var req request.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	input := dto.UpdateEventInput{
		Title:           req.Title,
		Description:     req.Description,
		Rules:           req.Rules,
		Recommendations: req.Recommendations,
		OrganizerNotes:  req.OrganizerNotes,
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		ClearDrawDate:   req.ClearDrawDate,
		EndDate:         req.EndDate,
		MaxParticipants: req.MaxParticipants,
		Budget:          req.Budget,
	}

	if err := h.uc.Update(r.Context(), id, userID, input); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.Delete(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) ActivateEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.Activate(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) FinishEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.Finish(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

