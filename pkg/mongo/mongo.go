package mongo

import (
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
}

func New(config Config) (*Mongo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Url))
	if err != nil {
		return nil, err
	}

	return &Mongo{
		client:  client,
		context: context.Background(),
	}, nil
}

func (m *Mongo) Close() {
	m.client.Disconnect(m.context)
}

func (m *Mongo) FindActiveConnection(UserID string, DeviceID string) (*Connection, error) {
	filter := bson.M{"user_id": UserID, "device_id": DeviceID}
	collection := m.client.Database("connections").Collection("connections")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection: %v", err)
	}

	var connection Connection
	if err = cursor.Decode(&connection); err != nil {
		return nil, fmt.Errorf("failed to decode connection: %v", err)
	}

	return &connection, nil
}

func (m *Mongo) InsertConnection(connection Connection) error {
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.InsertOne(context.Background(), connection)
	if err != nil {
		return fmt.Errorf("failed to insert connection: %v", err)
	}

	return nil
}

func (m *Mongo) DeleteConnection(userID string, deviceID string) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %v", err)
	}

	return nil
}

func (m *Mongo) HeartbeatConnection(userID string, deviceID string) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	update := bson.M{"$set": bson.M{"last_heartbeat": primitive.NewDateTimeFromTime(time.Now())}}
	collection := m.client.Database("connections").Collection("connections")
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update connection: %v", err)
	}

	return nil
}
