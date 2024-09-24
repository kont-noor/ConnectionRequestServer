package clientapp

import (
	"connection_request_server/internal/client"
	"time"

	"go.uber.org/zap"
)

func Run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Init client app")
	time.Sleep(4 * time.Second)

	client1 := client.New(client.Config{Host: "http://localhost:3000/api/v1", UserID: "15", DeviceID: "1", Log: logger})
	client2 := client.New(client.Config{Host: "http://localhost:3000/api/v1", UserID: "15", DeviceID: "2", Log: logger})
	client1.Connect()
	client2.Connect()

	for {
		time.Sleep(10 * time.Second)
		client1.Disconnect()
		client2.Connect()
		client1.Connect()
		time.Sleep(10 * time.Second)
		client2.Disconnect()
		client1.Connect()
		client2.Connect()
	}
}
