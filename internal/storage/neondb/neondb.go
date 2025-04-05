package neondb

import (
	"database/sql"
	"fmt"

	"github.com/kartikey1188/go-todo-list-v2/internal/config"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create students table: %w", err)
	}

	return &Postgres{
		Db: db,
	}, nil
}
