package mongo_adapter

import (
	"connection_request_server/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Url string
}

type Mongo struct {
	client  *mongo.Client
	context context.Context
	ttl     time.Duration
}

func New(config Config) (*Mongo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Url))
	if err != nil {
		return nil, err
	}

	return &Mongo{
		client:  client,
		context: context.Background(),
		ttl:     time.Duration(5000 * time.Millisecond),
	}, nil
}

func (m *Mongo) Close() {
	m.client.Disconnect(m.context)
}

func ConvertToMongoConnection(connection *domain.Connection) *Connection {
	return &Connection{
		ID:            primitive.NewObjectID(),
		UserID:        primitive.Symbol(connection.UserID),
		DeviceID:      primitive.Symbol(connection.DeviceID),
		LastHeartbeat: primitive.NewDateTimeFromTime(connection.LastHeartbeat),
		ConnectedAt:   primitive.NewDateTimeFromTime(connection.ConnectedAt),
	}
}

func ConvertToDomainConnection(connection *Connection) *domain.Connection {
	return &domain.Connection{
		ID:            fmt.Sprintf("%s", connection.ID),
		UserID:        fmt.Sprintf("%s", connection.UserID),
		DeviceID:      fmt.Sprintf("%s", connection.DeviceID),
		LastHeartbeat: connection.LastHeartbeat.Time(),
		ConnectedAt:   connection.ConnectedAt.Time(),
	}
}

func (m *Mongo) FindUserConnection(UserID string) (*domain.Connection, error) {
	filter := bson.M{"user_id": UserID, "last_heartbeat": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now().Add(-m.ttl))}}
	collection := m.client.Database("connections").Collection("connections")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection: %v", err)
	}

	var connection Connection
	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&connection); err != nil {
			return nil, fmt.Errorf("failed to decode connection: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no connection found for user ID %s", UserID)
	}

	//return &connection, nil
	return ConvertToDomainConnection(&connection), nil
}

func (m *Mongo) FindActiveConnection(UserID string, DeviceID string) (*domain.Connection, error) {
	filter := bson.M{"user_id": UserID, "device_id": DeviceID, "last_heartbeat": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now().Add(-m.ttl))}}
	collection := m.client.Database("connections").Collection("connections")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection: %v", err)
	}

	var connection Connection
	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&connection); err != nil {
			return nil, fmt.Errorf("failed to decode connection: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no connection found for user ID %s and device ID %s", UserID, DeviceID)
	}

	return ConvertToDomainConnection(&connection), nil
}

func (m *Mongo) InsertConnection(connection *domain.Connection) error {
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.InsertOne(context.Background(), ConvertToMongoConnection(connection))
	if err != nil {
		return fmt.Errorf("failed to insert connection: %v", err)
	}

	return nil
}

func (m *Mongo) DeleteConnection(userID string, deviceID string) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %v", err)
	}

	return nil
}

func (m *Mongo) HeartbeatConnection(userID string, deviceID string) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID, "last_heartbeat": bson.M{"$gt": primitive.NewDateTimeFromTime(time.Now().Add(-m.ttl))}}
	update := bson.M{"$set": bson.M{"last_heartbeat": primitive.NewDateTimeFromTime(time.Now())}}
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update connection: %v", err)
	}

	return nil
}
