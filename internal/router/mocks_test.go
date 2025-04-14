package router

import "net/http"

type mockHandlers struct {
	called map[string]bool
}

func newMockHandler() *mockHandlers {
	return &mockHandlers{
		called: make(map[string]bool),
	}
}

func (m *mockHandlers) Connect(w http.ResponseWriter, r *http.Request) {
	m.called["connected"] = true
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("connected"))
}

func (m *mockHandlers) Disconnect(w http.ResponseWriter, r *http.Request) {
	m.called["disconnected"] = true
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("disconnected"))
}

func (m *mockHandlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	m.called["heartbeat"] = true
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("heartbeat"))
}
