package clientapp

import (
	"connection_request_server/internal/client"
	"fmt"
	"time"
)

func Run() {
	fmt.Println("Init client app")
	time.Sleep(4 * time.Second)

	client := client.New(client.Config{Host: "http://localhost:3000", UserID: "1", DeviceID: "1"})
	client.Connect()

	for {
		client.Disconnect()
		time.Sleep(1 * time.Second)
		client.Connect()
		time.Sleep(1 * time.Second)
	}
}
