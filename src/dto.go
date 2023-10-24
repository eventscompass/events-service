// This file contains definitions of data transfer objects (DTOs) for the
// service API. HTTP requests containing a JSON body are decoded into a DTO
// before they are forwarded to the upstream (internal package). HTTP responses
// are in turn also stored in DTOs and are encoded as JSON strings before being
// returned to the caller.

package main

import (
	"time"
)

type EventDTO struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	StartDate time.Time     `json:"start_date"`
	EndDate   time.Time     `json:"end_date"`
	Location  LocationDTO   `json:"location"`
}

type LocationDTO struct {
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Country   string    `json:"country"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
	Halls     []HallDTO `json:"halls"`
}

type HallDTO struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
