package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/middleware"
	"secret-santa-backend/internal/usecase"
)

type ParticipantHandler struct {
	uc usecase.ParticipantUseCase
}

func NewParticipantHandler(uc usecase.ParticipantUseCase) *ParticipantHandler {
	return &ParticipantHandler{uc: uc}
}

func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	userID, err := middleware.GetUserID(r)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	participant, err := h.uc.Create(r.Context(), eventID, userID, entity.ParticipantRoleParticipant)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	resp := response.ParticipantResponse{
		ID:        participant.ID.String(),
		EventID:   participant.EventID.String(),
		UserID:    participant.UserID.String(),
		Role:      participant.Role,
		GiftSent:  participant.GiftSent,
		CreatedAt: participant.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ParticipantHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	participants, err := h.uc.GetByEvent(r.Context(), eventID)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	var resp []response.ParticipantResponse
	for _, p := range participants {
		resp = append(resp, response.ParticipantResponse{
			ID:        p.ID.String(),
			EventID:   p.EventID.String(),
			UserID:    p.UserID.String(),
			Role:      p.Role,
			GiftSent:  p.GiftSent,
			CreatedAt: p.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ParticipantHandler) MarkGiftSent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	err = h.uc.MarkGiftSent(r.Context(), id)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	err = h.uc.Delete(r.Context(), id)
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
