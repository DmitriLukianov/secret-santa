package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"
)

type UserHandler struct {
	uc      usecase.UserUseCase
	eventUC usecase.EventUseCase
}

func NewUserHandler(uc usecase.UserUseCase, eventUC usecase.EventUseCase) *UserHandler {
	return &UserHandler{
		uc:      uc,
		eventUC: eventUC,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	input := dto.CreateUserInput{
		Name:          req.Name,
		Email:         req.Email,
		OAuthID:       req.OAuthID,
		OAuthProvider: req.OAuthProvider,
	}

	_, err := h.uc.Create(r.Context(), input)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	user, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	resp := response.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAll(r.Context())
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}
	if err := helpers.ValidateStruct(&req); err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUserInput)
		return
	}

	input := dto.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	err = h.uc.Update(r.Context(), id, input)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteHTTPError(w, definitions.ErrInvalidUUID)
		return
	}

	err = h.uc.Delete(r.Context(), id)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetMyEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	events, err := h.eventUC.GetMyEvents(r.Context(), userID)
	if err != nil {
		response.WriteHTTPError(w, err)
		return
	}

	var resp []response.EventResponse
	for _, e := range events {
		resp = append(resp, response.EventResponse{
			ID:              e.ID.String(),
			Title:           e.Title,
			Description:     e.Description,
			Rules:           e.Rules,
			Recommendations: e.Recommendations,
			OrganizerID:     e.OrganizerID.String(),
			StartDate:       e.StartDate,
			DrawDate:        e.DrawDate,
			EndDate:         e.EndDate,
			Status:          string(e.Status),
			MaxParticipants: e.MaxParticipants,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
