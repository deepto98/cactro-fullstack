package db

import (
	"database/sql"
	"fmt"
)

// InitDB opens a PostgreSQL connection using the provided connection string
// and creates the necessary tables if they do not exist.
func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	// Verify the connection.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	// Create polls table.
	pollsTable := `
	CREATE TABLE IF NOT EXISTS polls (
		id SERIAL PRIMARY KEY,
		question TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(pollsTable); err != nil {
		return nil, fmt.Errorf("failed to create polls table: %w", err)
	}

	// Create poll_options table.
	optionsTable := `
	CREATE TABLE IF NOT EXISTS poll_options (
		id SERIAL PRIMARY KEY,
		poll_id INTEGER REFERENCES polls(id) ON DELETE CASCADE,
		option_text TEXT NOT NULL,
		vote_count INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(optionsTable); err != nil {
		return nil, fmt.Errorf("failed to create poll_options table: %w", err)
	}

	return db, nil
}
