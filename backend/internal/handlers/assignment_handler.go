package handlers

import (
	"encoding/json"
	"net/http"

	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/services"
)

type AssignmentHandler struct {
	service *services.AssignmentService
}

func NewAssignmentHandler(service *services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

func (h *AssignmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EventID    string `json:"event_id"`
		GiverID    string `json:"giver_id"`
		ReceiverID string `json:"receiver_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a := domain.Assignment{
		EventID:    req.EventID,
		GiverID:    req.GiverID,
		ReceiverID: req.ReceiverID,
	}

	err = h.service.Create(r.Context(), a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AssignmentHandler) GetMyAssignment(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	a, err := h.service.GetMy(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(a)
}
