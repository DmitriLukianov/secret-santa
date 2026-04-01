package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/middleware"
	"secret-santa-backend/internal/usecase"
)

type ChatHandler struct {
	uc usecase.ChatUseCase
}

func NewChatHandler(uc usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{uc: uc}
}

// GetRecipientChat — чат «Кому я Санта»
func (h *ChatHandler) GetRecipientChat(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
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

	messages, err := h.uc.GetRecipientChat(r.Context(), eventID, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.MessagesToResponse(messages))
}

// GetSenderChat — чат «Кто мой Санта»
func (h *ChatHandler) GetSenderChat(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
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

	messages, err := h.uc.GetSenderChat(r.Context(), eventID, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.MessagesToResponse(messages))
}

// SendMessage — отправить сообщение
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
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

	var req request.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	msg, err := h.uc.SendMessage(r.Context(), eventID, userID, req.Content)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.MessageToResponse(&msg))
}
