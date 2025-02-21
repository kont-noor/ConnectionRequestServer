package serverapp

import (
	adapter "connection_request_server/internal/mongo_adapter"
	"connection_request_server/internal/router"
	"connection_request_server/internal/service"
	"connection_request_server/pkg/server"
	"os"

	"go.uber.org/zap"
)

func Run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongoURL := os.Getenv("MONGO_URL")
	repositoryClient, err := adapter.New(adapter.Config{Url: mongoURL})
	if err != nil {
		logger.Sugar().Errorf("failed to create mongo client: %v, url: %s", err, mongoURL)
	}
	appService := service.New(service.Config{Repository: repositoryClient})
	appRouter := router.New(router.Config{APIHandlers: appService, Log: logger})

	serverHostname := os.Getenv("SERVER_HOSTNAME")
	serverPort := os.Getenv("SERVER_PORT")

	server := server.New(server.Config{Hostname: serverHostname, Port: serverPort, Router: appRouter, Log: logger})
	server.Run()
}
