package service

import "connection_request_server/internal/domain"

type MockRepository struct {
	Connections []*domain.Connection
	Error       error
}

func (t *MockRepository) FindUserConnection(userID string) (*domain.Connection, error) {
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

func (t *MockRepository) InsertConnection(connection *domain.Connection) error {
	t.Connections = append(t.Connections, connection)
	return nil
}

func (t *MockRepository) DeleteConnection(userID, deviceID string) error {
	return nil
}

func (t *MockRepository) HeartbeatConnection(userID, deviceID string) error {
	return nil
}
