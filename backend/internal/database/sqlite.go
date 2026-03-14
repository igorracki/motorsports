package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Manager struct {
	databaseConnection *sql.DB
}

func NewManager(databasePath string) (*Manager, error) {
	log.Printf("INFO: Initializing database at %s", databasePath)

	databaseConnection, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := databaseConnection.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	manager := &Manager{databaseConnection: databaseConnection}

	if err := manager.setup(); err != nil {
		return nil, fmt.Errorf("setting up database: %w", err)
	}

	log.Printf("INFO: Database initialized successfully")
	return manager, nil
}

func (manager *Manager) setup() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
	}

	for _, pragma := range pragmas {
		if _, err := manager.databaseConnection.Exec(pragma); err != nil {
			return fmt.Errorf("executing pragma %s: %w", pragma, err)
		}
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS profiles (
		user_id TEXT PRIMARY KEY,
		display_name TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS user_scores (
		user_id TEXT NOT NULL,
		score_type TEXT NOT NULL,
		season INTEGER,
		value INTEGER NOT NULL DEFAULT 0,
		updated_at TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id, score_type, season),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS predictions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		year INTEGER NOT NULL,
		round INTEGER NOT NULL,
		session_type TEXT NOT NULL,
		score INTEGER,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, year, round, session_type)
	);

	CREATE TABLE IF NOT EXISTS prediction_entries (
		prediction_id TEXT NOT NULL,
		position INTEGER NOT NULL,
		driver_id TEXT NOT NULL,
		PRIMARY KEY (prediction_id, position),
		FOREIGN KEY (prediction_id) REFERENCES predictions(id) ON DELETE CASCADE
	);`

	if _, err := manager.databaseConnection.Exec(schema); err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}

	return nil
}

func (manager *Manager) Close() error {
	log.Println("INFO: Closing database connection")
	return manager.databaseConnection.Close()
}

func (manager *Manager) DB() *sql.DB {
	return manager.databaseConnection
}
