package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/middleware"
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

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req request.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	participant, err := h.participantUC.GetByUserAndEvent(r.Context(), userID, uuid.MustParse(req.EventID))
	if err != nil {
		http.Error(w, "participant not found for this event", http.StatusBadRequest)
		return
	}

	wishlist, err := h.uc.Create(r.Context(), participant.ID, req.Visibility)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.WishlistResponse{
		ID:            wishlist.ID.String(),
		ParticipantID: wishlist.ParticipantID.String(),
		Visibility:    wishlist.Visibility,
		CreatedAt:     wishlist.CreatedAt,
		UpdatedAt:     wishlist.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *WishlistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "wishlistId")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		http.Error(w, "invalid wishlist id", http.StatusBadRequest)
		return
	}

	var req request.CreateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.WishlistItemResponse{
		ID:        item.ID.String(),
		Title:     item.Title,
		Link:      item.Link,
		ImageURL:  item.ImageURL,
		Comment:   item.Comment,
		CreatedAt: item.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *WishlistHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	eventIDStr := r.URL.Query().Get("eventId")
	if eventIDStr == "" {
		http.Error(w, "eventId query parameter is required (example: ?eventId=xxx)", http.StatusBadRequest)
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid eventId", http.StatusBadRequest)
		return
	}

	participant, err := h.participantUC.GetByUserAndEvent(r.Context(), userID, eventID)
	if err != nil {
		http.Error(w, "participant not found for this event", http.StatusNotFound)
		return
	}

	wishlist, err := h.uc.GetForUser(r.Context(), eventID, participant.ID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "you are not the santa") {
			http.Error(w, "you are not the santa for this participant", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.WishlistResponse{
		ID:            wishlist.ID.String(),
		ParticipantID: wishlist.ParticipantID.String(),
		Visibility:    wishlist.Visibility,
		CreatedAt:     wishlist.CreatedAt,
		UpdatedAt:     wishlist.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *WishlistHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "wishlistId")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		http.Error(w, "invalid wishlist id", http.StatusBadRequest)
		return
	}

	items, err := h.uc.GetItems(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.WishlistItemResponse
	for _, item := range items {
		resp = append(resp, response.WishlistItemResponse{
			ID:        item.ID.String(),
			Title:     item.Title,
			Link:      item.Link,
			ImageURL:  item.ImageURL,
			Comment:   item.Comment,
			CreatedAt: item.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
