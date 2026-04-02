package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Manager struct {
	db *sql.DB
}

func NewManager(databasePath string) (*Manager, error) {
	slog.Info("Initializing database", "path", databasePath)

	dir := filepath.Dir(databasePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	manager := &Manager{db: db}

	if err := manager.setup(); err != nil {
		return nil, fmt.Errorf("setting up database: %w", err)
	}

	if err := manager.migrate(); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	slog.Info("Database initialized and migrated successfully")
	return manager, nil
}

func (manager *Manager) setup() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
	}

	for _, pragma := range pragmas {
		if _, err := manager.db.Exec(pragma); err != nil {
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

	CREATE TABLE IF NOT EXISTS predictions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		year INTEGER NOT NULL,
		round INTEGER NOT NULL,
		session_type TEXT NOT NULL,
		score INTEGER,
		revalidate_until TIMESTAMP,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, year, round, session_type)
	);

	CREATE TABLE IF NOT EXISTS prediction_entries (
		prediction_id TEXT NOT NULL,
		position INTEGER NOT NULL,
		driver_id TEXT NOT NULL,
		correct INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (prediction_id, position),
		FOREIGN KEY (prediction_id) REFERENCES predictions(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS friend_requests (
		id TEXT PRIMARY KEY,
		sender_id TEXT NOT NULL,
		receiver_id TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(sender_id, receiver_id)
	);

	CREATE TABLE IF NOT EXISTS friendships (
		user_id TEXT NOT NULL,
		friend_id TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id, friend_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := manager.db.Exec(schema); err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}

	return nil
}

func (manager *Manager) migrate() error {
	if _, err := manager.db.Exec("DROP TABLE IF EXISTS user_scores"); err != nil {
		return fmt.Errorf("dropping user_scores table: %w", err)
	}
	return nil
}

func (manager *Manager) Close() error {
	slog.Info("Closing database connection")
	return manager.db.Close()
}

func (manager *Manager) DB() *sql.DB {
	return manager.db
}

func (manager *Manager) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	transaction, err := manager.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		}
	}()

	if err := fn(transaction); err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
