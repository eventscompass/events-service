# Events-Service

The `Events` micro-service manages the events on the platform.
It can be used to create new events or to retrieve existing events.
Events can be retrieved using their unique ID, or by their name.

## REST API
| method | route                           | description                   |
|--------|---------------------------------|-------------------------------|
|  GET   | `/api/events/id/<uid>`          | retrieve an event by its ID   |
|  GET   | `/api/events/name/<event_name>` | retrieve an event by its name |
|  GET   | `/api/events`                   | retrieve all events           |
|  POST  | `/api/events`                   | create a new event            |