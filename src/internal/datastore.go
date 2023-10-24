package internal

import (
	"context"
	"time"
)

// EventsContainer abstracts the database layer for storing events.
type EventsContainer interface {

	// Create creates a new event in the container. This method
	// returns the ID associated with the created event.
	Create(_ context.Context, e Event) (string, error)

	// GetByID retrieves the event with the given id from the container.
	GetByID(_ context.Context, id string) (Event, error)

	// GetByName retrieves the event with the given name from the container.
	GetByName(_ context.Context, name string) (Event, error)

	// GetAll retrieves all events from the container.
	GetAll(_ context.Context) ([]Event, error)
}

// Event represents a single document entry in the container.
type Event struct {
	ID        string
	Name      string
	Duration  time.Duration
	StartDate time.Time
	EndDate   time.Time
	Location  Location
}

// Location represents a real location with opening and closing time.
type Location struct {
	Name      string
	Address   string
	Country   string
	OpenTime  time.Time
	CloseTime time.Time
	Halls     []Hall
}

// Hall is the room where the event will be taking place.
type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
