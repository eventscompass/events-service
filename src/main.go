package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/events-service/src/internal/mongodb"
	"github.com/eventscompass/service-framework/pubsub"
	"github.com/eventscompass/service-framework/pubsub/rabbitmq"
	"github.com/eventscompass/service-framework/service"
)

// EventsService manages the events. It can be used to create new
// events or retrieve existing ones.
type EventsService struct {
	service.BaseService

	// restHandler is used to start an http server, that serves
	// http requests to the rest api of the service.
	restHandler http.Handler

	// eventBus is used for publishing and subscribing to messages.
	eventsBus service.MessageBus

	// eventsDB is used to read and store elements in a container database.
	eventsDB internal.EventsContainer

	// cfg is used to configure the service.
	cfg *Config
}

// Init implements the [service.CloudService] interface.
func (s *EventsService) Init(ctx context.Context) error {
	// Parse the env variables.
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("%w: env parse: %v", service.ErrUnexpected, err)
	}
	s.cfg = &cfg

	// Init the database layer.
	mongoCfg := mongodb.Config(s.cfg.EventsDB)
	db, err := mongodb.NewMongoDBContainer(ctx, &mongoCfg)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	s.eventsDB = db

	// Init the message bus,
	busCfg := rabbitmq.Config(s.cfg.EventsMQ)
	bus, err := rabbitmq.NewAMQPBus(&busCfg, pubsub.EventsExchange)
	if err != nil {
		return fmt.Errorf("init mq: %w", err)
	}
	s.eventsBus = bus

	// Init the rest API of the service.
	s.initREST()

	return nil
}

func main() {
	service.Start(&EventsService{})
}
