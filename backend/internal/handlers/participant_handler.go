package handlers

import (
	"encoding/json"
	"net/http"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type ParticipantHandler struct {
	service *services.ParticipantService
}

func NewParticipantHandler(service *services.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{service: service}
}

func (h *ParticipantHandler) JoinEvent(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "event_id")

	var req struct {
		UserID string `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := domain.Participant{
		EventID: eventID,
		UserID:  req.UserID,
	}

	err = h.service.JoinEvent(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "joined event",
	})
}
func (h *ParticipantHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "event_id")

	participants, err := h.service.GetParticipants(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(participants)
}

func (h *ParticipantHandler) LeaveEvent(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "event_id")
	userID := chi.URLParam(r, "user_id")

	err := h.service.LeaveEvent(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
