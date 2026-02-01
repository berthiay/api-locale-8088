package httpserver

import (
	"context"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const roleKey contextKey = "role"

var tokenRoles = map[string]string{
	os.Getenv("BEZA_ADMIN_TOKEN"):  "admin",
	os.Getenv("BEZA_VIEWER_TOKEN"): "viewer",
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		role, ok := tokenRoles[token]
		if !ok || role == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
