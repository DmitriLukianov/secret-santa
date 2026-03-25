package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/controller/http/v1/request"
	"secret-santa-backend/internal/controller/http/v1/response"
	"secret-santa-backend/internal/dto"
	eventusecase "secret-santa-backend/internal/usecase/event"
)

type EventHandler struct {
	uc *eventusecase.UseCase
}

func NewEventHandler(uc *eventusecase.UseCase) *EventHandler {
	return &EventHandler{uc: uc}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req request.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.CreateEventInput{
		Name:            req.Name,
		Description:     req.Description,
		Rules:           req.Rules,
		Recommendations: req.Recommendations,
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		EndDate:         req.EndDate,
		MaxParticipants: req.MaxParticipants,
	}

	organizerID := r.Header.Get("X-Organizer-ID")
	created, err := h.uc.Create(r.Context(), input, organizerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response.EventResponse{
		ID:              created.ID.String(),
		Name:            created.Name,
		Description:     created.Description,
		OrganizerID:     created.OrganizerID.String(),
		StartDate:       created.StartDate,
		DrawDate:        created.DrawDate,
		EndDate:         created.EndDate,
		Rules:           created.Rules,
		Recommendations: created.Recommendations,
		Status:          created.Status,
		MaxParticipants: created.MaxParticipants,
		CreatedAt:       created.CreatedAt,
		UpdatedAt:       created.UpdatedAt,
	})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response.EventResponse{
		ID:              event.ID.String(),
		Name:            event.Name,
		Description:     event.Description,
		OrganizerID:     event.OrganizerID.String(),
		StartDate:       event.StartDate,
		DrawDate:        event.DrawDate,
		EndDate:         event.EndDate,
		Rules:           event.Rules,
		Recommendations: event.Recommendations,
		Status:          event.Status,
		MaxParticipants: event.MaxParticipants,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
	})
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.uc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]response.EventResponse, 0, len(events))
	for _, e := range events {
		resp = append(resp, response.EventResponse{
			ID:              e.ID.String(),
			Name:            e.Name,
			Description:     e.Description,
			OrganizerID:     e.OrganizerID.String(),
			StartDate:       e.StartDate,
			DrawDate:        e.DrawDate,
			EndDate:         e.EndDate,
			Rules:           e.Rules,
			Recommendations: e.Recommendations,
			Status:          e.Status,
			MaxParticipants: e.MaxParticipants,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req request.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.UpdateEventInput{
		Name:            req.Name,
		Description:     req.Description,
		Rules:           req.Rules,
		Recommendations: req.Recommendations,
		StartDate:       req.StartDate,
		DrawDate:        req.DrawDate,
		EndDate:         req.EndDate,
		Status:          req.Status,
		MaxParticipants: req.MaxParticipants,
	}

	updated, err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response.EventResponse{
		ID:              updated.ID.String(),
		Name:            updated.Name,
		Description:     updated.Description,
		OrganizerID:     updated.OrganizerID.String(),
		StartDate:       updated.StartDate,
		DrawDate:        updated.DrawDate,
		EndDate:         updated.EndDate,
		Rules:           updated.Rules,
		Recommendations: updated.Recommendations,
		Status:          updated.Status,
		MaxParticipants: updated.MaxParticipants,
		CreatedAt:       updated.CreatedAt,
		UpdatedAt:       updated.UpdatedAt,
	})
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.uc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
