package internal

import (
	"context"
	"time"
)

// EventsContainer abstracts the database layer for storing events.
type EventsContainer interface {

	// Create creates a new entry in the given collection in the
	// container. This method returns the ID associated with the
	// created entry.
	Create(_ context.Context, e Event) (string, error)

	// GetByID retrieves the entry with the given id from the
	// given collection in the container. This function returns
	// [service.ErrNotFound] if the requested item is not in the
	// container.
	GetByID(_ context.Context, id string) (Event, error)

	// GetByID retrieves the entry with the given name from the
	// given collection in the container. This function returns
	// [service.ErrNotFound] if the requested item is not in the
	// container.
	GetByName(_ context.Context, name string) (Event, error)

	// GetAll retrieves all entries from the given collection
	// from the container.
	GetAll(_ context.Context) ([]Event, error)
}

// Event represents an event entry in the container.
type Event struct {
	ID        string
	Name      string
	Duration  time.Duration
	StartDate time.Time
	EndDate   time.Time
	Location  Location
}

// Location represents a location entry in the container.
type Location struct {
	ID        string
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

var (
	// Collections.
	EventsCollection    string = "events"
	LocationsCollection string = "locations"
)
