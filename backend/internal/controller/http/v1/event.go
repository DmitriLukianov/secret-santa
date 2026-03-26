package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/middleware"
	"secret-santa-backend/internal/usecase"
)

type EventHandler struct {
	uc usecase.EventUseCase
}

func NewEventHandler(uc usecase.EventUseCase) *EventHandler {
	return &EventHandler{uc: uc}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req request.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.EventResponse{
		ID:              event.ID.String(),
		Title:           event.Title,
		Description:     event.Description,
		Rules:           event.Rules,
		Recommendations: event.Recommendations,
		OrganizerID:     event.OrganizerID.String(),
		StartDate:       event.StartDate,
		DrawDate:        event.DrawDate,
		EndDate:         event.EndDate,
		Status:          string(event.Status),
		MaxParticipants: event.MaxParticipants,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	event, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := response.EventResponse{
		ID:              event.ID.String(),
		Title:           event.Title,
		Description:     event.Description,
		Rules:           event.Rules,
		Recommendations: event.Recommendations,
		OrganizerID:     event.OrganizerID.String(),
		StartDate:       event.StartDate,
		DrawDate:        event.DrawDate,
		EndDate:         event.EndDate,
		Status:          string(event.Status),
		MaxParticipants: event.MaxParticipants,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.uc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.EventResponse
	for _, e := range events {
		resp = append(resp, response.EventResponse{
			ID:              e.ID.String(),
			Title:           e.Title,
			Description:     e.Description,
			Rules:           e.Rules,
			Recommendations: e.Recommendations,
			OrganizerID:     e.OrganizerID.String(),
			StartDate:       e.StartDate,
			DrawDate:        e.DrawDate,
			EndDate:         e.EndDate,
			Status:          string(e.Status),
			MaxParticipants: e.MaxParticipants,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req request.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		Status:          req.Status,
		MaxParticipants: req.MaxParticipants,
	}

	err = h.uc.Update(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	err = h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) FinishEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	err = h.uc.Finish(r.Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "organizer") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Событие успешно завершено",
	})
}
