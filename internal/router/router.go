package router

import "net/http"

const (
	apiPath        = "/api/v1"
	connectPath    = "/connect"
	disconnectPath = "/disconnect"
	heartbeatPath  = "/heartbeat"
)

type APIHandlers interface {
	Connect(w http.ResponseWriter, r *http.Request)
	Disconnect(w http.ResponseWriter, r *http.Request)
	Heartbeat(w http.ResponseWriter, r *http.Request)
}

func New(apiHandlers APIHandlers) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST "+apiPath+connectPath, apiHandlers.Connect)
	router.HandleFunc("POST "+apiPath+disconnectPath, apiHandlers.Disconnect)
	router.HandleFunc("POST "+apiPath+heartbeatPath, apiHandlers.Heartbeat)
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connection request server"))
	})

	return router
}
