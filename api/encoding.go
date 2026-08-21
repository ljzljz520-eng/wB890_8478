package api

import (
	"encoding/json"
	"memorialstation/domain"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error(), Code: code})
}
func decodeRecord(r *http.Request) (*domain.Record, error) {
	var record domain.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		return nil, err
	}
	if err := domain.NormalizeRecord(&record); err != nil {
		return nil, err
	}
	return &record, nil
}
func methodAllowed(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}
func requestActor(r *http.Request) string {
	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "anonymous"
	}
	return actor
}
func queryValue(r *http.Request, key string) string { return r.URL.Query().Get(key) }
