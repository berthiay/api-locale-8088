package httpserver

import "database/sql"

var db *sql.DB

func SetDB(d *sql.DB) {
	db = d
}

func getDB() *sql.DB {
	return db
}
