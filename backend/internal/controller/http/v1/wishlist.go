package v1

import (
	"encoding/json"
	"net/http"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/usecase/wishlist"

	"github.com/go-chi/chi/v5"
)

type WishlistHandler struct {
	uc *wishlist.UseCase
}

func NewWishlistHandler(uc *wishlist.UseCase) *WishlistHandler {
	return &WishlistHandler{uc: uc}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req request.WishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.CreateWishlistInput{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Link:        req.Link,
		ImageURL:    req.ImageURL,
		Visibility:  req.Visibility,
	}

	err := h.uc.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *WishlistHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	items, err := h.uc.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.WishlistResponse

	for _, wItem := range items {
		resp = append(resp, response.WishlistResponse{
			ID:          wItem.ID,
			UserID:      wItem.UserID,
			Title:       wItem.Title,
			Description: wItem.Description,
			Link:        wItem.Link,
			ImageURL:    wItem.ImageURL,
			Visibility:  wItem.Visibility,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
