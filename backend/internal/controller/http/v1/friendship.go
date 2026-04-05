package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"
)

type FriendshipHandler struct {
	uc usecase.FriendshipUseCase
}

func NewFriendshipHandler(uc usecase.FriendshipUseCase) *FriendshipHandler {
	return &FriendshipHandler{uc: uc}
}

func (h *FriendshipHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var req request.SendFriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	addresseeID, err := uuid.Parse(req.AddresseeID)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	friendship, err := h.uc.SendRequest(r.Context(), userID, addresseeID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.FriendshipToResponse(&friendship))
}

func (h *FriendshipHandler) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.AcceptRequest(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendshipHandler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.DeclineRequest(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendshipHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.RemoveFriend(r.Context(), id, userID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendshipHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	friends, err := h.uc.GetFriends(r.Context(), userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.FriendshipsToResponse(friends))
}

func (h *FriendshipHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	requests, err := h.uc.GetPendingRequests(r.Context(), userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.FriendshipsToResponse(requests))
}
