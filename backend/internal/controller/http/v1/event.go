package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"secret-santa-backend/internal/controller/http/v1/request" // ← ВОТ ЭТО ДОБАВЬ
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

// CREATE EVENT
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req request.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.CreateEventInput{
		Name:        req.Name,
		Description: req.Description,
		OrganizerID: req.OrganizerID,
		StartDate:   req.StartDate,
		DrawDate:    req.DrawDate,
		EndDate:     req.EndDate,
	}

	err := h.uc.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GET EVENT BY ID
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := response.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		OrganizerID: event.OrganizerID,
		StartDate:   event.StartDate,
		DrawDate:    event.DrawDate,
		EndDate:     event.EndDate,
	}

	json.NewEncoder(w).Encode(resp)
}

// GET ALL EVENTS
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.uc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []response.EventResponse

	for _, e := range events {
		resp = append(resp, response.EventResponse{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			OrganizerID: e.OrganizerID,
			StartDate:   e.StartDate,
			DrawDate:    e.DrawDate,
			EndDate:     e.EndDate,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

// UPDATE EVENT
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req request.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := dto.UpdateEventInput{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DELETE EVENT
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
