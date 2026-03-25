package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/usecase"
)

type AssignmentHandler struct {
	uc usecase.AssignmentUseCase
}

func NewAssignmentHandler(uc usecase.AssignmentUseCase) *AssignmentHandler {
	return &AssignmentHandler{uc: uc}
}

// Draw — запускает жеребьёвку
func (h *AssignmentHandler) Draw(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	err = h.uc.Draw(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Жеребьёвка успешно проведена",
	})
}

// GetByEvent — показывает результаты жеребьёвки
func (h *AssignmentHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	assignments, err := h.uc.GetByEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.AssignmentResponse
	for _, a := range assignments {
		resp = append(resp, response.AssignmentResponse{
			ID:         a.ID.String(),
			EventID:    a.EventID.String(),
			GiverID:    a.GiverID.String(),
			ReceiverID: a.ReceiverID.String(),
			CreatedAt:  a.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
