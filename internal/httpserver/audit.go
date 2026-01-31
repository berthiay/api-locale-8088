package httpserver

import (
	"log"
	"net/http"
	"time"
)

// Audit consigne chaque requête (méthode, chemin, rôle, durée).
func Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		role, _ := r.Context().Value(roleKey).(string)
		log.Printf("audit: method=%s path=%s role=%s duration=%s", r.Method, r.URL.Path, role, time.Since(start))
	})
}
