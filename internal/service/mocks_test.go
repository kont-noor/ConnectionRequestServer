package service

import "connection_request_server/internal/domain"

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
	return nil
}

func (t *mockRepository) HeartbeatConnection(userID, deviceID string) error {
	return nil
}
