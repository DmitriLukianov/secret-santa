package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response" // ← ВОТ ЭТО
	"secret-santa-backend/internal/dto"
	participantusecase "secret-santa-backend/internal/usecase/participant"
)

type ParticipantHandler struct {
	uc *participantusecase.UseCase
}

func NewParticipantHandler(uc *participantusecase.UseCase) *ParticipantHandler {
	return &ParticipantHandler{uc: uc}
}

// ADD PARTICIPANT
func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	var req request.AddParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.AddParticipantInput{
		EventID: eventID,
		UserID:  req.UserID,
	}

	err := h.uc.Add(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GET PARTICIPANTS BY EVENT
func (h *ParticipantHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	participants, err := h.uc.GetByEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.ParticipantResponse

	for _, p := range participants {
		resp = append(resp, response.ParticipantResponse{
			ID:      p.ID,
			EventID: p.EventID,
			UserID:  p.UserID,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

// DELETE PARTICIPANT
func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
