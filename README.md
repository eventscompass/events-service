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


## Configuration
The service is configured using environment variables.

| name                            | default  | description                                                     |
|---------------------------------|----------|-----------------------------------------------------------------|
| HTTP_SERVER_LISTEN              | :8080    | The address for the service to listen on for http requests.     |
| HTTP_SERVER_READ_HEADER_TIMEOUT | 10s      | How long to wait for reading the http request headers.          |
| HTTP_SERVER_READ_TIMEOUT        | 10s      | How long to wait for reading the http requests, including body. |
| HTTP_SERVER_WRITE_TIMEOUT       | 30s      | How long to wait to process requests and generate a response.   |
| HTTP_SERVER_DUMP_REQUESTS       |          |                                                                 |
| MESSAGE_BUS_HOST                |          | The host url for connecting to a message bus.                   |
| MESSAGE_BUS_PORT                |          | The port on which the message bus listens.                      |
| MESSAGE_BUS_USERNAME            |          | The username for connecting to the message bus.                 |
| MESSAGE_BUS_PASSWORD            |          | The password for connecting to the message bus.                 |
| EVENTS_MONGO_HOST               |          | The host url for connecting to a MongoDB server.                |
| EVENTS_MONGO_PORT               |          | The port on which the database server listens.                  |
| EVENTS_MONGO_USERNAME           |          | The username for connecting to the server.                      |
| EVENTS_MONGO_PASSWORD           |          | The password for connecting to the server.                      |
| EVENTS_MONGO_DATABASE           |          | The name of the database that is allocated for this service.    |
