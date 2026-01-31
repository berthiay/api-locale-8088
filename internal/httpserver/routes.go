package httpserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/berthiay/api-locale-8088/internal/storage"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/health", health)

	// reports: viewer + admin
	mux.Handle("/v1/reports", Auth(http.HandlerFunc(reports)))

	// admin only
	mux.Handle("/v1/decision-proposals", Auth(RequireRole("admin", http.HandlerFunc(decisionProposals))))
	mux.Handle("/v1/validation-requests", Auth(RequireRole("admin", http.HandlerFunc(validationRequests))))
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func reports(w http.ResponseWriter, r *http.Request) {
	handleListCreate(w, r, "reports")
}

func decisionProposals(w http.ResponseWriter, r *http.Request) {
	handleListCreate(w, r, "decision_proposals")
}

func validationRequests(w http.ResponseWriter, r *http.Request) {
	handleListCreate(w, r, "validation_requests")
}

func handleListCreate(w http.ResponseWriter, r *http.Request, table string) {
	d := getDB()
	if d == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := storage.List(d, table)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(items)
		return

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req, err := storage.DecodeCreateRequest(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		created, err := storage.Create(d, table, req.Content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(created)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
