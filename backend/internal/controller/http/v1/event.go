package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"
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

	input := dto.CreateEventInput{
		Title:           req.Title,
		Description:     req.Description,
		Rules:           req.Rules,
		Recommendations: req.Recommendations,
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		EndDate:         req.EndDate,
		MaxParticipants: req.MaxParticipants,
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
	events, err := h.uc.GetAll(r.Context())
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.EventsToResponse(events))
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
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
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		EndDate:         req.EndDate,
		MaxParticipants: req.MaxParticipants,
	}

	if err := h.uc.Update(r.Context(), id, input); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) OpenInvitation(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, h.uc.OpenInvitation)
}

func (h *EventHandler) CloseRegistration(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, h.uc.CloseRegistration)
}

func (h *EventHandler) StartDrawing(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, h.uc.StartDrawing)
}

func (h *EventHandler) FinishEvent(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, h.uc.Finish)
}

func (h *EventHandler) CancelEvent(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, h.uc.Cancel)
}

func (h *EventHandler) changeStatus(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, id, userID uuid.UUID) error,
) {
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

	if err := action(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Статус события успешно изменён"})
}
