# Events-Service

[![CircleCI](https://dl.circleci.com/status-badge/img/circleci/RGExmu1KKSYDZZz3vWH7qN/JozX5aRwsFCZY23aXBiZCb.svg?style=svg&circle-token=b51a25bc4fc74c08f2d33f1764a0380083b374ae)](https://dl.circleci.com/status-badge/redirect/circleci/RGExmu1KKSYDZZz3vWH7qN/JozX5aRwsFCZY23aXBiZCb)

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