package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berthiay/api-locale-8088/internal/httpserver"
)

func main() {
	addr := env("BEZA_ADDR", ":8088")

	mux := http.NewServeMux()
	httpserver.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Bezalel OS API listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
