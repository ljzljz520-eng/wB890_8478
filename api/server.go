package api

import (
	"encoding/json"
	"memorialstation/archive"
	"memorialstation/domain"
	"memorialstation/importx"
	"memorialstation/review"
	"memorialstation/search"
	"memorialstation/storage"
	"memorialstation/workflow"
	"net/http"
	"strings"
)

type Server struct {
	store    *storage.Store
	workflow *workflow.Service
	review   *review.Service
	archive  *archive.Service
	search   *search.Service
	importer *importx.Service
}

func New(store *storage.Store) *Server {
	return &Server{store: store, workflow: workflow.New(store), review: review.New(store), archive: archive.New(store), search: search.New(store), importer: importx.New(store)}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/batches/", s.batch)
	mux.HandleFunc("/search", s.query)
	return mux
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func readJSON(r *http.Request, target any) error { return json.NewDecoder(r.Body).Decode(target) }
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		records, err := s.store.ListRecords()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, records)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var record domain.Record
	if err := readJSON(r, &record); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if err := s.workflow.CreateRecord(&record, "api"); err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, record)
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeJSON(w, 400, map[string]string{"error": "id required"})
		return
	}
	if r.Method == http.MethodGet {
		item, err := s.store.GetRecord(id)
		if err != nil {
			writeJSON(w, 404, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, item)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	action := r.URL.Query().Get("action")
	var err error
	switch action {
	case "submit":
		err = s.workflow.Submit(id, "api", "api")
	case "confirm":
		err = s.workflow.Confirm(id, "api", "api")
	case "archive":
		_, err = s.archive.Archive(id, "api", "api")
	case "publish":
		err = s.workflow.Publish(id, "api")
	default:
		err = domain.ErrInvalidInput
	}
	if err != nil {
		writeJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	item, err := s.store.GetRecord(id)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, item)
}

func (s *Server) batch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/batches/")
	if id == "" {
		writeJSON(w, 400, map[string]string{"error": "batch id required"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		summary, err := s.archive.Summary(id)
		if err != nil {
			writeJSON(w, 404, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, summary)
	case http.MethodPost:
		var batch domain.Batch
		if err := readJSON(r, &batch); err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		batch.ID = id
		if err := s.workflow.StartBatch(&batch, "api"); err != nil {
			writeJSON(w, 422, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 201, batch)
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) query(w http.ResponseWriter, r *http.Request) {
	query := domain.SearchQuery{BatchID: r.URL.Query().Get("batch"), Status: r.URL.Query().Get("status"), Visibility: r.URL.Query().Get("visibility"), StudentName: r.URL.Query().Get("student"), Tag: r.URL.Query().Get("tag")}
	result, err := s.search.Query(query)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, result)
}
