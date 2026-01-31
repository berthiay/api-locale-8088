package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berthiay/api-locale-8088/internal/httpserver"
	"github.com/berthiay/api-locale-8088/internal/storage"
)

func main() {
	addr := env("BEZA_ADDR", ":8088")

	// Base SQLite locale + schema
	db, err := storage.InitDB("bezalel.db", "001_init.sql")
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer db.Close()

	httpserver.SetDB(db)

	mux := http.NewServeMux()
	httpserver.RegisterRoutes(mux)

	// audit global
	handler := httpserver.Audit(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
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
