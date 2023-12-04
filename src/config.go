package main

// Config encapsulates the configuration of the service.
type Config struct {

	// EventsDB encapsulates the configuration of the database
	// layer used by the service.
	EventsDB DBConfig

	// BusConfig encapsulates the configuration for the message
	// bus used by the service.
	EventsMQ BusConfig
}

// DBConfig encapsulates the configuration of the database layer
// used by the service.
type DBConfig struct {
	Host     string `env:"MONGO_DB_HOST" envDefault:"mongodb"`
	Port     int    `env:"MONGO_DB_PORT" envDefault:"27017"`
	Username string `env:"MONGO_DB_USERNAME" envDefault:"user"`
	Password string `env:"MONGO_DB_PASSWORD" envDefault:"password"`
	Database string `env:"MONGO_DB_DATABASE" envDefault:"events"`
}

// BusConfig encapsulates the configuration for the message bus
// used by the service.
type BusConfig struct {
	Host     string `env:"RABBIT_MQ_HOST"`
	Port     int    `env:"RABBIT_MQ_PORT"`
	Username string `env:"RABBIT_MQ_USERNAME"`
	Password string `env:"RABBIT_MQ_PASSWORD"`
}
