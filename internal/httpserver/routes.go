package httpserver

import (
	"encoding/json"
	"net/http"
)

// Enregistre toutes les routes.
func RegisterRoutes(mux *http.ServeMux) {
	// /v1/health reste public
	mux.HandleFunc("/v1/health", health)

	// /v1/reports : accessible aux rôles admin et viewer
	mux.Handle("/v1/reports", Auth(http.HandlerFunc(reports)))

	// /v1/decision-proposals : uniquement admin
	mux.Handle("/v1/decision-proposals", Auth(RequireRole("admin", http.HandlerFunc(decisionProposals))))

	// /v1/validation-requests : uniquement admin
	mux.Handle("/v1/validation-requests", Auth(RequireRole("admin", http.HandlerFunc(validationRequests))))
}

// /v1/health renvoie le statut OK
func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// /v1/reports : liste fictive de rapports
func reports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode([]string{"report1", "report2"})
}

// /v1/decision-proposals : liste fictive de propositions
func decisionProposals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode([]string{"proposal1", "proposal2"})
}

// /v1/validation-requests : liste fictive de demandes de validation
func validationRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode([]string{"validation1", "validation2"})
}
