package main

import (
	"github.com/eventscompass/service-framework/service"
)

// Bus implements the [service.CloudService] interface.
func (s *EventsService) Bus() service.MessageBus {
	return s.eventsBus
}
