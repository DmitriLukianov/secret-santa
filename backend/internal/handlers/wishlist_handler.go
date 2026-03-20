package handlers

import (
	"encoding/json"
	"net/http"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type WishlistHandler struct {
	service *services.WishlistService
}

func NewWishlistHandler(service *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{service: service}
}

func (h *WishlistHandler) CreateWishlist(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "event_id")

	var req struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wishlist := domain.Wishlist{
		EventID: eventID,
		UserID:  req.UserID,
		Text:    req.Text,
	}

	err = h.service.Create(r.Context(), wishlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "wishlist created",
	})
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "event_id")
	userID := chi.URLParam(r, "user_id")

	wishlist, err := h.service.GetByUser(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wishlist)
}

func (h *WishlistHandler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var req struct {
		Text *string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), id, req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
