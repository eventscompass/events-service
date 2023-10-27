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
func (s *EventsService) initREST() error {
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
	return nil
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
	var e internal.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the event.
	id, err := h.eventsDB.Create(ctx, e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish to the message queue.
	payload := &pubsub.EventCreated{
		ID:         e.ID,
		Name:       e.Name,
		LocationID: e.Location.Name,
		Start:      e.StartDate,
		End:        e.EndDate,
	}
	if err := h.eventsBus.Publish(ctx, payload); err != nil {
		// log the error.. no idea how to handle it, but obviously it should be handled!
		log.Println("failed to publish message:", *payload, err)
	}

	// Write the response.
	w.Header().Set("location", fmt.Sprintf("%s/id/%s", r.URL.Path, id))
	w.WriteHeader(http.StatusCreated)
}

func (h *restHandler) readByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	id := chi.URLParam(r, "id")

	event, err := h.eventsDB.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *restHandler) readByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request key.
	name := chi.URLParam(r, "name")

	event, err := h.eventsDB.GetByName(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&event)
}

func (h *restHandler) readAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	events, err := h.eventsDB.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	_ = json.NewEncoder(w).Encode(&events)
}
