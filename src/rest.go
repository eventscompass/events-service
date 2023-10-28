package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/service-framework/pubsub"
	"github.com/eventscompass/service-framework/service"
)

// REST implements the [service.CloudService] interface.
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
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	err := h.eventsDB.Create(ctx, internal.EventsCollection, event)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Publish to the message queue.
	payload := &pubsub.EventCreated{
		ID:         event.ID,
		Name:       event.Name,
		LocationID: event.Location.ID,
		Start:      event.StartDate,
		End:        event.EndDate,
	}
	if err := h.eventsBus.Publish(ctx, payload); err != nil {
		// TODO: handle this error somehow. Check this out:
		// https://cloud.google.com/pubsub/docs/samples/pubsub-publish-with-error-handler
		log.Println("failed to publish message:", *payload, err)
	}

	// Write the response.
	w.Header().Set("location", fmt.Sprintf("%s/id/%s", r.URL.Path, event.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h *restHandler) readByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	id := chi.URLParam(r, "id")

	// Get the event.
	event, err := h.eventsDB.GetByID(ctx, internal.EventsCollection, id)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *restHandler) readByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	name := chi.URLParam(r, "name")

	// Get the event.
	event, err := h.eventsDB.GetByName(ctx, internal.EventsCollection, name)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *restHandler) readAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all events.
	events, err := h.eventsDB.GetAll(ctx, internal.EventsCollection)
	if err != nil {
		service.HTTPError(ctx, w, err)
		return
	}

	// Write the response.
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&events)
}
