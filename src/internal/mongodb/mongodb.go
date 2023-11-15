package mongodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	. "github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/service-framework/service"
)

// Config holds configuration variables for connecting to a Mongo database.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// MongoDBContainer is a container backed by a Mongo database.
//
//nolint:revive // this type will probably not be used as mongodb.MongoDBContainer
type MongoDBContainer struct {
	client   *mongo.Client
	database *mongo.Database
}

var (
	_ EventsContainer = (*MongoDBContainer)(nil)
	_ io.Closer       = (*MongoDBContainer)(nil)
)

// NewMongoDBContainer creates a new [MongoDBContainer] instance.
func NewMongoDBContainer(ctx context.Context, cfg *Config) (*MongoDBContainer, error) {
	connInfo := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(connInfo)
	clientOptions.SetAuth(options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	})

	// Use once.Do to make sure that in this service we have only one mongodb
	// connection, even if this function is called multiple times.
	var err error
	once.Do(func() { client, err = mongo.Connect(ctx, clientOptions) })
	if err != nil {
		return nil, service.Unexpected(ctx, fmt.Errorf("mongo connect: %w", err))
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx) //nolint:errcheck // intentional
		return nil, service.Unexpected(ctx, fmt.Errorf("ping mongo: %w", err))
	}

	database := client.Database(cfg.Database)
	return &MongoDBContainer{
		client:   client,
		database: database,
	}, nil
}

// Create implements the [EventsContainer] interface.
func (m *MongoDBContainer) Create(
	ctx context.Context,
	collection string,
	data any,
) error {
	c := m.database.Collection(collection)
	if _, err := c.InsertOne(ctx, data); err != nil {
		return service.Unexpected(ctx, fmt.Errorf("insert one: %w", err))
	}
	return nil
}

// GetByID implements the [EventsContainer] interface.
func (m *MongoDBContainer) GetByID(
	ctx context.Context,
	collection string,
	id string,
) (any, error) {
	return m.findOne(ctx, collection, "id", id)
}

// GetByName implements the [EventsContainer] interface.
func (m *MongoDBContainer) GetByName(
	ctx context.Context,
	collection string,
	name string,
) (any, error) {
	return m.findOne(ctx, collection, "name", name)
}

// GetAll implements the [EventsContainer] interface.
func (m *MongoDBContainer) GetAll(
	ctx context.Context,
	collection string,
) ([]any, error) {
	// Get all elements from the requested collection.
	c := m.database.Collection(collection)
	cursor, err := c.Find(ctx, bson.D{})
	if err != nil {
		return nil, service.Unexpected(ctx, fmt.Errorf("find: %w", err))
	}

	// Use context.Background() to ensure Close completes even if the ctx passed
	// to this function has errored.
	defer cursor.Close(context.Background()) //nolint:errcheck, contextcheck // intentional

	res := make([]any, 0)

	// Infer the type of the elements and decode them.
	switch collection {
	case EventsCollection:
		var elems []Event
		if err := cursor.All(ctx, &elems); err != nil {
			return nil, service.Unexpected(ctx, fmt.Errorf("cursor all: %w", err))
		}
		for _, e := range elems {
			res = append(res, e)
		}
	case LocationsCollection:
		var elems []Location
		if err := cursor.All(ctx, &elems); err != nil {
			return nil, service.Unexpected(ctx, fmt.Errorf("cursor all: %w", err))
		}
		for _, e := range elems {
			res = append(res, e)
		}
	default:
		return nil, fmt.Errorf(
			"%w: unknown collection %q", service.ErrNotAllowed, collection)
	}

	return res, nil
}

func (m *MongoDBContainer) findOne(
	ctx context.Context,
	collection string,
	filterKey string,
	filterValue any,
) (any, error) {
	// Get the element from the collection.
	c := m.database.Collection(collection)
	one := c.FindOne(ctx, bson.M{filterKey: filterValue})
	if err := one.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: %v", service.ErrNotFound, err)
		}
		return nil, service.Unexpected(ctx, fmt.Errorf("find one: %w", err))
	}

	// Infer the type of the requested element and decode it.
	switch collection {
	case EventsCollection:
		var elem Event
		if err := one.Decode(&elem); err != nil {
			return nil, service.Unexpected(ctx, fmt.Errorf("decode one: %w", err))
		}
		return elem, nil
	case LocationsCollection:
		var elem Location
		if err := one.Decode(&elem); err != nil {
			return nil, service.Unexpected(ctx, fmt.Errorf("decode one: %w", err))
		}
		return elem, nil
	default:
		return nil, fmt.Errorf(
			"%w: unknown collection %q", service.ErrNotAllowed, collection)
	}
}

// Close implements the [io.Closer] interface.
func (m *MongoDBContainer) Close() error {
	// Disconnect the client by waiting up to 10 seconds for
	// in-progress operations to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:gomnd // intentional
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		return service.Unexpected(ctx, err)
	}
	return nil
}

var (
	// Use a singleton to make sure only one connection is open.
	once   sync.Once
	client *mongo.Client
)
