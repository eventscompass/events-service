package main

import (
	"github.com/eventscompass/service-framework/service"
)

// Config encapsulates the configuration of the service.
type Config struct {
	// EventsDB encapsulates the configuration of the database
	// layer used by the service.
	EventsDB DBConfig

	// BusConfig encapsulates the configuration for the message
	// bus used by the service.
	EventsMQ service.BusConfig
}

// DBConfig encapsulates the configuration of the database layer
// used by the service.
type DBConfig struct {
	Host     string `env:"EVENTS_MONGO_HOST" envDefault:"mongodb"`
	Port     int    `env:"EVENTS_MONGO_PORT" envDefault:"27017"`
	Username string `env:"EVENTS_MONGO_USERNAME" envDefault:"user"`
	Password string `env:"EVENTS_MONGO_PASSWORD" envDefault:"password"`
	Database string `env:"EVENTS_MONGO_DATABASE" envDefault:"events"`
}
