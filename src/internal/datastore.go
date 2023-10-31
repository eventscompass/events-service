package internal

import (
	"context"
	"io"
	"time"
)

// EventsContainer abstracts the database layer for storing events.
type EventsContainer interface {
	io.Closer

	// Create creates a new entry in the given collection in the
	// container.
	Create(_ context.Context, collection string, data any) error

	// GetByID retrieves the entry with the given id from the
	// given collection in the container. This function returns
	// [service.ErrNotFound] if the requested item is not in the
	// container. This function returns [service.ErrNotAllowed]
	// if the requested collection is not in the container.
	GetByID(_ context.Context, collection string, id string) (any, error)

	// GetByID retrieves the entry with the given name from the
	// given collection in the container. This function returns
	// [service.ErrNotFound] if the requested item is not in the
	// container. This function returns [service.ErrNotAllowed]
	// if the requested collection is not in the container.
	GetByName(_ context.Context, collection string, name string) (any, error)

	// GetAll retrieves all entries from the given collection
	// from the container. This function returns
	// [service.ErrNotAllowed] if the requested collection is not
	// in the container.
	GetAll(_ context.Context, collection string) ([]any, error)
}

// Event represents an event entry in the container.
type Event struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	StartDate time.Time     `json:"start_date"`
	EndDate   time.Time     `json:"end_date"`
	Location  Location      `json:"location"`
}

// Location represents a location entry in the container.
type Location struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Country   string    `json:"country"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
	Halls     []Hall    `json:"halls"`
}

// Hall is the room where the event will be taking place.
type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity int    `json:"capacity"`
}

var (
	// Collections.
	EventsCollection    string = "events"
	LocationsCollection string = "locations"
)
