package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func createRequest(method, path string, params *requestParams) *http.Request {
	var buf bytes.Buffer

	if params != nil {
		if err := json.NewEncoder(&buf).Encode(params); err != nil {
			panic("Failed to encode params: " + err.Error())
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}
