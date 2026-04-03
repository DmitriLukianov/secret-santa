package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"
)

type InvitationHandler struct {
	uc usecase.InvitationUseCase
}

func NewInvitationHandler(uc usecase.InvitationUseCase) *InvitationHandler {
	return &InvitationHandler{uc: uc}
}

func (h *InvitationHandler) GenerateInvite(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var req request.CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	expiresIn := 7 * 24 * time.Hour
	if req.ExpiresIn != "" {
		if d, err := time.ParseDuration(req.ExpiresIn); err == nil {
			expiresIn = d
		} else {
			response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
			return
		}
	}

	input := dto.CreateInvitationInput{
		EventID:   req.EventID,
		ExpiresIn: expiresIn,
	}

	resp, err := h.uc.GenerateInvite(r.Context(), input, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *InvitationHandler) JoinByInvite(w http.ResponseWriter, r *http.Request) {
	var req request.JoinByInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrUnauthorized)
		return
	}

	input := dto.JoinByInvitationInput{
		Token:  req.Token,
		UserID: userID,
	}

	if err := h.uc.JoinByInvite(r.Context(), input); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Вы успешно присоединились к событию",
	})
}
