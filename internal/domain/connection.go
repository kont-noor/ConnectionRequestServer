package domain

import (
	"time"
)

type Connection struct {
	ID            string
	UserID        string
	DeviceID      string
	ConnectedAt   time.Time
	LastHeartbeat time.Time
}
