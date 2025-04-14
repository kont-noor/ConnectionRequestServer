package service

import (
	"connection_request_server/internal/domain"
	"errors"
	"time"
)

type mockRepository struct {
	Connections []*domain.Connection
	Error       error
}

func (t *mockRepository) FindUserConnection(userID string) (*domain.Connection, error) {
	if t.Error != nil {
		return nil, t.Error
	}

	for _, connection := range t.Connections {
		if connection.UserID == userID {
			return connection, nil
		}
	}

	return nil, nil
}

func (t *mockRepository) InsertConnection(connection *domain.Connection) error {
	t.Connections = append(t.Connections, connection)
	return nil
}

func (t *mockRepository) DeleteConnection(userID, deviceID string) error {
	if t.Error != nil {
		return t.Error
	}

	for i, connection := range t.Connections {
		if connection.UserID == userID && connection.DeviceID == deviceID {
			t.deleteConnAt(i)
			return nil
		}
	}

	return errors.New("connection not found")
}

func (t *mockRepository) HeartbeatConnection(userID, deviceID string) error {
	if t.Error != nil {
		return t.Error
	}

	for _, connection := range t.Connections {
		if connection.UserID == userID && connection.DeviceID == deviceID {
			connection.LastHeartbeat = time.Now()
			return nil
		}
	}

	return errors.New("connection not found")
}

func (t *mockRepository) deleteConnAt(index int) {
	if len(t.Connections) == 1 && index == 0 {
		t.Connections = []*domain.Connection{}
		return
	}

	t.Connections[index] = t.Connections[len(t.Connections)-1]
	t.Connections = t.Connections[:len(t.Connections)-1]
}
