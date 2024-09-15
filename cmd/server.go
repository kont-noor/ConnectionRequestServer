package main

import (
	"connection_request_server/pkg/server"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	server := server.New(server.Config{Port: "3000"})
	server.Run()
}
