package main

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/events-service/src/internal/mongodb"
	"github.com/eventscompass/service-framework/service"
)

// EventsService manages the events. It can be used to create new
// events or retrieve existing ones.
type EventsService struct {
	service.BaseService

	eventsDB internal.EventsContainer
	cfg      *Config
}

// Init implements the [CloudService] interface.
func (s *EventsService) Init(ctx context.Context) error {
	// Parse the env variables.
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	s.cfg = &cfg

	mongoCfg := mongodb.Config(s.cfg.EventsDB)
	db, err := mongodb.NewMongoDBContainer(ctx, &mongoCfg)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	s.eventsDB = db

	if err := s.initREST(); err != nil {
		return fmt.Errorf("init rest: %w", err)
	}
	return nil
}

func main() {
	service.Start(&EventsService{})
}
