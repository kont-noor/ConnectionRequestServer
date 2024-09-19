package clientapp

import (
	"connection_request_server/internal/client"
	"fmt"
	"time"
)

func Run() {
	fmt.Println("Init client app")
	time.Sleep(4 * time.Second)

	client1 := client.New(client.Config{Host: "http://localhost:3000", UserID: "10", DeviceID: "1"})
	client2 := client.New(client.Config{Host: "http://localhost:3000", UserID: "10", DeviceID: "2"})
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
