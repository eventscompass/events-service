package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	. "github.com/eventscompass/events-service/src/internal"
	"github.com/eventscompass/service-framework/service"
)

// Config holds configuration variables for connecting to a Mongo database.
type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	Collection string
}

// MongoDBContainer is a container backed by a Mongo database.
type MongoDBContainer struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

var _ EventsContainer = (*MongoDBContainer)(nil)

// NewMongoDBContainer creates a new [MongoDBContainer] instance.
func NewMongoDBContainer(ctx context.Context, cfg *Config) (*MongoDBContainer, error) {
	connInfo := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(connInfo)
	clientOptions.SetAuth(options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	database := client.Database(cfg.Database)
	collection := database.Collection(cfg.Collection)
	return &MongoDBContainer{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// Create creates a new entry in the database collection.
func (m *MongoDBContainer) Create(ctx context.Context, e Event) (string, error) {
	_, err := m.collection.InsertOne(ctx, e)
	if err != nil {
		return "", service.Unexpected(ctx, err)
	}
	return e.ID, nil
}

// GetByID returns an entry from the database collection with the given ID.
func (m *MongoDBContainer) GetByID(ctx context.Context, id string) (Event, error) {
	return m.findOne(ctx, "id", id)
}

// GetByName returns an entry from the database collection with the given name.
func (m *MongoDBContainer) GetByName(ctx context.Context, name string) (Event, error) {
	return m.findOne(ctx, "name", name)
}

// GetAll returns all entries in the database collection.
func (m *MongoDBContainer) GetAll(ctx context.Context) ([]Event, error) {
	cursor, err := m.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, service.Unexpected(ctx, err)
	}
	defer cursor.Close(ctx)

	events := make([]Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, service.Unexpected(ctx, err)
	}
	return events, nil
}

func (m *MongoDBContainer) findOne(
	ctx context.Context,
	filterKey string,
	filterValue any,
) (Event, error) {
	one := m.collection.FindOne(ctx, bson.M{filterKey: filterValue})
	if err := one.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Event{}, fmt.Errorf("%w: %v", service.ErrNotFound, err)
		}
		return Event{}, service.Unexpected(ctx, err)
	}

	var e Event
	if err := one.Decode(&e); err != nil {
		return Event{}, service.Unexpected(ctx, err)
	}
	return e, nil
}
