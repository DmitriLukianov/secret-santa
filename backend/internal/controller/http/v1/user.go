package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/usecase" // ← публичный интерфейс UserUseCase
)

type UserHandler struct {
	uc usecase.UserUseCase // ← теперь используем интерфейс из contracts.go
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// CreateUser — теперь передаём все обязательные поля (OAuthID + OAuthProvider)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.CreateUserInput{
		Name:          req.Name,
		Email:         req.Email,
		OAuthID:       req.OAuthID,       // ← добавлено
		OAuthProvider: req.OAuthProvider, // ← добавлено
	}

	_, err := h.uc.Create(r.Context(), input) // ← теперь возвращает entity.User
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUserByID — теперь используем GetByID + uuid
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := response.UserResponse{
		ID:    user.ID.String(), // ← uuid → string для JSON
		Name:  user.Name,
		Email: user.Email,
	}

	json.NewEncoder(w).Encode(resp)
}

// GetUsers
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.UserResponse
	for _, u := range users {
		resp = append(resp, response.UserResponse{
			ID:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

// UpdateUser
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	err = h.uc.Update(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
