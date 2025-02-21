package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.Symbol   `bson:"user_id"`
	DeviceID      primitive.Symbol   `bson:"device_id"`
	ConnectedAt   primitive.DateTime `bson:"connected_at"`
	LastHeartbeat primitive.DateTime `bson:"last_heartbeat"`
}
