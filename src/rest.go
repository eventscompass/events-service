package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/service-framework/pubsub"
	"github.com/eventscompass/service-framework/service"
)

func (s *EventsService) REST() http.Handler {
	return s.restHandler
}

// initREST initializes the handler for the rest server part of the service.
// This function creates a router and registers with that router the handlers
// for the http endpoints.
func (s *EventsService) initREST() {
	restHandler := &restHandler{
		eventsDB:  s.eventsDB,
		eventsBus: s.eventsBus,
	}
	mux := chi.NewMux()

	// API routes.
	mux.Get("/api/events/id/{id}", restHandler.readByID)
	mux.Get("/api/events/name/{name}", restHandler.readByName)
	mux.Get("/api/events", restHandler.readAll)
	mux.Post("/api/events", restHandler.create)

	// Health check.
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I am healthy and strong, buddy!")
	}))

	s.restHandler = mux
}

// restHandler handles http requests. It is the bridge between the rest api and
// the business logic. Every rest endpoint exposed by the server will be served
// by calling one of the handler methods.
type restHandler struct {
	eventsDB  internal.EventsContainer
	eventsBus service.MessageBus
}

func (h *restHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body.
	var event internal.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		service.HTTPError(ctx, w, fmt.Errorf("%w: %v", service.ErrBadRequest, err))
		return
	}

	// Create the event.
	slog.Info("request to create event", slog.Any("event", event))
	err := h.eventsDB.Create(ctx, internal.EventsCollection, event)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}
	slog.Info("event successfully created")

	// Publish to the message queue.
	payload := pubsub.EventCreated{
		ID:         event.ID,
		Name:       event.Name,
		LocationID: event.Location.ID,
		Start:      event.StartDate,
		End:        event.EndDate,
	}
	topic := pubsub.EventCreatedTopic
	body, err := json.Marshal(&payload)
	if err != nil {
		slog.Error("failed to marshal for publishing", err)
	}
	if err == nil {
		if pErr := h.eventsBus.Publish(ctx, topic, body); pErr != nil {
			slog.Error("failed to publish", slog.String("topic", topic), pErr)
		}
		slog.Info("publish message", slog.String("topic", topic), slog.Any("message", payload))
	}

	// Write the response.
	w.Header().Set("Location", fmt.Sprintf("%s/id/%s", r.URL.Path, event.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h *restHandler) readByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	id := chi.URLParam(r, "id")

	// Get the event.
	slog.Info("request to read event", slog.String("id", id))
	event, err := h.eventsDB.GetByID(ctx, internal.EventsCollection, id)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(&event); err != nil {
		slog.Info("failed to write response", err)
	}
}

func (h *restHandler) readByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	name := chi.URLParam(r, "name")

	// Get the event.
	slog.Info("request to read event", slog.String("name", name))
	event, err := h.eventsDB.GetByName(ctx, internal.EventsCollection, name)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(&event); err != nil {
		slog.Info("failed to write response", err)
	}
}

func (h *restHandler) readAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all events.
	slog.Info("request to read all events")
	events, err := h.eventsDB.GetAll(ctx, internal.EventsCollection)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(&events); err != nil {
		slog.Info("failed to write response", err)
	}
}
