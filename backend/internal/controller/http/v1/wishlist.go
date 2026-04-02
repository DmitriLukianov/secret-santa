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

type WishlistHandler struct {
	uc            usecase.WishlistUseCase
	participantUC usecase.ParticipantUseCase
}

func NewWishlistHandler(uc usecase.WishlistUseCase, participantUC usecase.ParticipantUseCase) *WishlistHandler {
	return &WishlistHandler{
		uc:            uc,
		participantUC: participantUC,
	}
}

// Create — создание вишлиста
func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var req request.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	participant, err := h.participantUC.GetByUserAndEvent(r.Context(), userID, eventID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	wishlist, err := h.uc.Create(r.Context(), participant.ID, req.Visibility)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.WishlistToResponse(&wishlist))
}

// GetByParticipant — главный метод для Санты
func (h *WishlistHandler) GetByParticipant(w http.ResponseWriter, r *http.Request) {
	participantIDStr := chi.URLParam(r, "participantId")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	eventIDStr := r.URL.Query().Get("eventId")
	if eventIDStr == "" {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	wishlist, err := h.uc.GetForUser(r.Context(), eventID, participantID, userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.WishlistToResponse(wishlist))
}

func (h *WishlistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "wishlistId")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	var req request.CreateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	item, err := h.uc.AddItem(
		r.Context(),
		wishlistID,
		req.Title,
		&req.Link,
		&req.ImageURL,
		&req.Comment,
	)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.WishlistItemToResponse(&item))
}

// GetByUser — теперь правильно получает СВОЙ вишлист (для владельца)
func (h *WishlistHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	eventIDStr := r.URL.Query().Get("eventId")
	if eventIDStr == "" {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	participant, err := h.participantUC.GetByUserAndEvent(r.Context(), userID, eventID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	// 🔥 КРИТИЧЕСКИЙ ФИКС: владелец смотрит свой вишлист напрямую
	wishlist, err := h.uc.GetByParticipant(r.Context(), participant.ID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.WishlistToResponse(wishlist))
}

func (h *WishlistHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "wishlistId")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	items, err := h.uc.GetItems(r.Context(), wishlistID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var resp []response.WishlistItemResponse
	for _, item := range items {
		resp = append(resp, response.WishlistItemToResponse(&item))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateItem — обновление товара
// UpdateItem — обновление товара
func (h *WishlistHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	var req request.UpdateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	item, err := h.uc.UpdateItem(r.Context(), itemID, req.Title, &req.Link, &req.ImageURL, &req.Comment)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.WishlistItemToResponse(&item))
}

// DeleteItem — удаление товара
func (h *WishlistHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	if err := h.uc.DeleteItem(r.Context(), itemID); err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
