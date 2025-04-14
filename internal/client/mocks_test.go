package client

import (
	"net/http/httptest"

	"net/http"
	"sync"
)

type mockServer struct {
	HeartbeatCount int
	mu             *sync.Mutex
	Server         *httptest.Server
}

func newMockServer(connectResponse, disconnectResponse int) *mockServer {
	mockS := mockServer{
		mu: &sync.Mutex{},
	}
	mockS.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/connect":
			w.WriteHeader(connectResponse)
		case "/disconnect":
			w.WriteHeader(disconnectResponse)
		case "/heartbeat":
			mockS.mu.Lock()
			mockS.HeartbeatCount++
			mockS.mu.Unlock()
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return &mockS
}
