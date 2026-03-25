package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/entity" // ← добавили для роли
	"secret-santa-backend/internal/usecase"
)

type ParticipantHandler struct {
	uc usecase.ParticipantUseCase
}

func NewParticipantHandler(uc usecase.ParticipantUseCase) *ParticipantHandler {
	return &ParticipantHandler{uc: uc}
}

// Add
func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	// TODO: позже брать из JWT middleware
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	participant, err := h.uc.Create(r.Context(), eventID, userID, entity.ParticipantRoleParticipant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// GetByEvent
func (h *ParticipantHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	participants, err := h.uc.GetByEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// MarkGiftSent
func (h *ParticipantHandler) MarkGiftSent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid participant id", http.StatusBadRequest)
		return
	}

	err = h.uc.MarkGiftSent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete
func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid participant id", http.StatusBadRequest)
		return
	}

	err = h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
