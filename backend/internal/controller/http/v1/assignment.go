package v1

import (
	"encoding/json"
	"net/http"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/usecase/assignment"

	"github.com/go-chi/chi/v5"
)

type AssignmentHandler struct {
	uc *assignment.UseCase
}

func NewAssignmentHandler(uc *assignment.UseCase) *AssignmentHandler {
	return &AssignmentHandler{uc: uc}
}

// 🎁 POST /events/{eventId}/assign
func (h *AssignmentHandler) Draw(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	input := dto.GenerateAssignmentInput{
		EventID: eventID,
	}

	err := h.uc.Draw(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// 📦 GET /events/{eventId}/assignments
func (h *AssignmentHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	assignments, err := h.uc.GetByEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.AssignmentResponse

	for _, a := range assignments {
		resp = append(resp, response.AssignmentResponse{
			ID:         a.ID,
			EventID:    a.EventID,
			GiverID:    a.GiverID,
			ReceiverID: a.ReceiverID,
		})
	}

	json.NewEncoder(w).Encode(resp)
}
