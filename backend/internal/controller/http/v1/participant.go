package v1

import (
	"encoding/json"
	"net/http"

	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ParticipantHandler struct {
	uc      usecase.ParticipantUseCase
	eventUC usecase.EventUseCase
}

func NewParticipantHandler(uc usecase.ParticipantUseCase, eventUC usecase.EventUseCase) *ParticipantHandler {
	return &ParticipantHandler{uc: uc, eventUC: eventUC}
}


func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
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

	participant, err := h.uc.Create(r.Context(), eventID, userID, definitions.ParticipantRoleParticipant)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.ParticipantToResponse(&participant))
}

func (h *ParticipantHandler) GetByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	requesterID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	// Проверяем права доступа: должен быть участником или организатором
	event, err := h.eventUC.GetByID(r.Context(), eventID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	if event.OrganizerID != requesterID {
		// Проверяем членство через прямой запрос (не через пагинацию)
		if _, err := h.uc.GetByUserAndEvent(r.Context(), requesterID, eventID); err != nil {
			response.WriteHTTPError(w, definitions.ErrForbidden)
			return
		}
	}

	pg := helpers.ParsePagination(r)
	participants, total, err := h.uc.GetByEventPaged(r.Context(), eventID, pg.Limit, pg.Offset())
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(helpers.NewPagedResponse(response.ParticipantsToResponse(participants), total, pg))
}

func (h *ParticipantHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	participant, err := h.uc.GetByUserAndEvent(r.Context(), userID, eventID)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrParticipantNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.ParticipantToResponse(participant))
}


func (h *ParticipantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	requesterID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	if err := h.uc.Delete(r.Context(), id, requesterID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
