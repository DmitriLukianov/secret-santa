package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"
)

type AssignmentHandler struct {
	uc usecase.AssignmentUseCase
}

func NewAssignmentHandler(uc usecase.AssignmentUseCase) *AssignmentHandler {
	return &AssignmentHandler{uc: uc}
}

func (h *AssignmentHandler) Draw(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.Draw(r.Context(), eventID, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Жеребьёвка успешно проведена",
	})
}

func (h *AssignmentHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	assignments, err := h.uc.GetByEvent(r.Context(), eventID, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

	json.NewEncoder(w).Encode(resp)
}
