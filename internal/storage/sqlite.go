package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
}

func InitDB(dbPath string, schemaPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(schemaBytes)); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func List(db *sql.DB, table string) ([]Item, error) {
	q := "SELECT id, created_at, content FROM " + table + " ORDER BY id DESC"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// IMPORTANT: liste vide = [] (pas null)
	out := make([]Item, 0)

	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.CreatedAt, &it.Content); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func Create(db *sql.DB, table string, content string) (Item, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return Item{}, errors.New("content is required")
	}

	res, err := db.Exec("INSERT INTO "+table+" (content) VALUES (?)", content)
	if err != nil {
		return Item{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Item{}, err
	}

	var it Item
	err = db.QueryRow("SELECT id, created_at, content FROM "+table+" WHERE id = ?", id).
		Scan(&it.ID, &it.CreatedAt, &it.Content)
	if err != nil {
		return Item{}, err
	}
	return it, nil
}

type CreateRequest struct {
	Content string `json:"content"`
}

func DecodeCreateRequest(body []byte) (CreateRequest, error) {
	var req CreateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return CreateRequest{}, err
	}
	return req, nil
}
