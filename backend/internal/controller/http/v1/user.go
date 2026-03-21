package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	userusecase "secret-santa-backend/internal/usecase/user"
)

type UserHandler struct {
	uc *userusecase.UseCase
}

func NewUserHandler(uc *userusecase.UseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	err := h.uc.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.UserResponse

	for _, u := range users {
		resp = append(resp, response.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
