package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/events-service/src/internal/mongodb"
	"github.com/eventscompass/service-framework/service"
)

// EventsService manages the events. It can be used to create new
// events or retrieve existing ones.
type EventsService struct {
	service.BaseService

	// restHandler is used to start an http server, that serves
	// http requests to the rest api of the service.
	restHandler http.Handler

	// eventsDB is used to read and store events in a container database.
	eventsDB internal.EventsContainer

	// cfg is used to configure the service.
	cfg *Config
}

// Init implements the [CloudService] interface.
func (s *EventsService) Init(ctx context.Context) error {
	// Parse the env variables.
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	s.cfg = &cfg

	// Init the database layer.
	mongoCfg := mongodb.Config(s.cfg.EventsDB)
	db, err := mongodb.NewMongoDBContainer(ctx, &mongoCfg)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	s.eventsDB = db

	// Init the rest API of the service.
	if err := s.initREST(); err != nil {
		return fmt.Errorf("init rest: %w", err)
	}
	return nil
}

func main() {
	service.Start(&EventsService{})
}