package main

type Config struct {
	// EventsDB encapsulates the configuration of the database
	// layer used by the service.
	EventsDB DBConfig
}

// DBConfig encapsulates the configuration of the database layer
// used by the service.
type DBConfig struct {
	Host       string `env:"EVENTS_MONGO_HOST" envDefault:"mongodb"`
	Port       int    `env:"EVENTS_MONGO_PORT" envDefault:"27017"`
	Username   string `env:"EVENTS_MONGO_USERNAME" envDefault:"user"`
	Password   string `env:"EVENTS_MONGO_PASSWORD" envDefault:"password"`
	Database   string `env:"EVENTS_MONGO_DATABASE" envDefault:"events"`
	Collection string `env:"EVENTS_MONGO_COLLECTION" envDefault:"events"`
}