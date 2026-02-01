package httpserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/berthiay/api-locale-8088/internal/storage"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/health", health)

	// OpenAPI (public)
	mux.HandleFunc("/openapi.yaml", openapiYAML)

	// reports: viewer + admin
	mux.Handle("/v1/reports", Auth(http.HandlerFunc(reports)))

	// admin only
	mux.Handle("/v1/decision-proposals", Auth(RequireRole("admin", http.HandlerFunc(decisionProposals))))
	mux.Handle("/v1/validation-requests", Auth(RequireRole("admin", http.HandlerFunc(validationRequests))))
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
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
		writeErr(w, http.StatusInternalServerError, "db not initialized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := storage.List(d, table)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, items)
		return

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "cannot read body")
			return
		}
		req, err := storage.DecodeCreateRequest(body)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json body (expected: {\"content\":\"...\"})")
			return
		}
		created, err := storage.Create(d, table, req.Content)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, created)
		return

	default:
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
