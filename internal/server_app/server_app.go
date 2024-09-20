package serverapp

import (
	"connection_request_server/internal/router"
	"connection_request_server/internal/service"
	"connection_request_server/pkg/mongo"
	"connection_request_server/pkg/server"
	"log"
)

func Run() {
	mongoConfig := mongo.Config{Url: "mongodb://mongo:mongo@localhost:27017"}
	mongoClient, err := mongo.New(mongoConfig)
	if err != nil {
		log.Fatalf("failed to create mongo client: %v", err)
	}

	appService := service.New(service.Config{Mongo: mongoClient})
	appRouter := router.New(appService)

	server := server.New(server.Config{Port: "3000", Router: appRouter})
	server.Run()
}
