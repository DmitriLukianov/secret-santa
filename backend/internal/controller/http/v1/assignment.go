package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/middleware"
	"secret-santa-backend/internal/usecase"
)

type AssignmentHandler struct {
	uc usecase.AssignmentUseCase
}

func NewAssignmentHandler(uc usecase.AssignmentUseCase) *AssignmentHandler {
	return &AssignmentHandler{uc: uc}
}

func (h *AssignmentHandler) Draw(w http.ResponseWriter, r *http.Request) {

	userID, err := middleware.GetUserID(r)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	err = h.uc.Draw(r.Context(), eventID, userID)
	if err != nil {
		if errors.Is(err, definitions.ErrNotOrganizer) {
			writeHTTPError(w, definitions.ErrForbidden)
			return
		}
		writeHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Жеребьёвка успешно проведена",
	})
}

func (h *AssignmentHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	assignments, err := h.uc.GetByEvent(r.Context(), eventID, userID)
	if err != nil {
		writeHTTPError(w, err)
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
