package httpserver

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const roleKey contextKey = "role"

// Jetons et rôles associés.
// Remplace éventuellement les valeurs admin123 et viewer123 par tes propres jetons.
var tokenRoles = map[string]string{
	"admin123":  "admin",
	"viewer123": "viewer",
}

// Auth vérifie le jeton Bearer et stocke le rôle dans le contexte.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		role, ok := tokenRoles[token]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole autorise l’accès uniquement si le rôle correspond.
func RequireRole(required string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(roleKey).(string)
		if role != required {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
