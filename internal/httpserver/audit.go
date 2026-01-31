package httpserver

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Audit consigne chaque requête (méthode, chemin, rôle, durée).
func Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		role := guessRoleFromAuthHeader(r.Header.Get("Authorization"))
		log.Printf("audit: method=%s path=%s role=%s duration=%s", r.Method, r.URL.Path, role, time.Since(start))
	})
}

func guessRoleFromAuthHeader(authHeader string) string {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if role, ok := tokenRoles[token]; ok {
		return role
	}
	return ""
}
