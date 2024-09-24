package serverapp

import (
	"connection_request_server/internal/router"
	"connection_request_server/internal/service"
	"connection_request_server/pkg/mongo"
	"connection_request_server/pkg/server"

	"go.uber.org/zap"
)

func Run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongoClient, err := mongo.New(mongo.Config{Url: "mongodb://mongo:mongo@localhost:27017"})
	if err != nil {
		logger.Sugar().Errorf("failed to create mongo client: %v", err)
	}
	appService := service.New(service.Config{Repository: mongoClient})
	appRouter := router.New(router.Config{APIHandlers: appService, Log: logger})

	server := server.New(server.Config{Port: "3000", Router: appRouter, Log: logger})
	server.Run()
}
