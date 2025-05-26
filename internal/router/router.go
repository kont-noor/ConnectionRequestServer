package router

import (
	"bytes"
	"connection_request_server/internal/metrics"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

const (
	apiPath        = "/api/v1"
	connectPath    = "/connect"
	disconnectPath = "/disconnect"
	heartbeatPath  = "/heartbeat"
	healthPath     = "/health"
)

type APIHandlers interface {
	Connect(w http.ResponseWriter, r *http.Request)
	Disconnect(w http.ResponseWriter, r *http.Request)
	Heartbeat(w http.ResponseWriter, r *http.Request)
}

type Config struct {
	APIHandlers       APIHandlers
	MetricsMiddleware *metrics.Middleware
	Log               *zap.Logger
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		logger.Info("Received request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("params", string(bodyBytes)),
		)

		rw := &ResponseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Info("Sent response",
			zap.Int("status_code", rw.statusCode),
			zap.String("response_body", rw.body.String()),
		)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func New(config Config) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle(apiPath+connectPath, applyMiddlewares(config, config.APIHandlers.Connect))
	router.Handle(apiPath+disconnectPath, applyMiddlewares(config, config.APIHandlers.Disconnect))
	router.Handle(apiPath+heartbeatPath, applyMiddlewares(config, config.APIHandlers.Heartbeat))

	router.Handle(healthPath, http.HandlerFunc(healthCheckHandler))

	config.MetricsMiddleware.Routes(router)

	return router
}

func applyMiddlewares(config Config, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return config.MetricsMiddleware.Middleware(LoggingMiddleware(config.Log, http.HandlerFunc(handler)))
}
