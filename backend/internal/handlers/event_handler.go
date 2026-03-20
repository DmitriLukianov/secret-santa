package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OrganizerID string `json:"organizer_id"`
		StartDate   string `json:"start_date"`
		DrawDate    string `json:"draw_date"`
		EndDate     string `json:"end_date"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse(time.RFC3339, req.StartDate)
	drawDate, _ := time.Parse(time.RFC3339, req.DrawDate)
	endDate, _ := time.Parse(time.RFC3339, req.EndDate)

	event := domain.Event{
		Name:        req.Name,
		Description: req.Description,
		OrganizerID: req.OrganizerID,
		StartDate:   startDate,
		DrawDate:    drawDate,
		EndDate:     endDate,
	}

	err = h.eventService.CreateEvent(r.Context(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "event created",
	})
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {

	events, err := h.eventService.GetEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	event, err := h.eventService.GetEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.eventService.UpdateEvent(r.Context(), id, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := h.eventService.DeleteEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
