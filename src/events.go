package main

import (
	"github.com/eventscompass/service-framework/service"
)

func (s *EventsService) Bus() service.MessageBus {
	return s.eventsBus
}
