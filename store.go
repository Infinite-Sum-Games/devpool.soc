package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var localDb *sql.DB

type Entry struct {
	ID        int64
	StreamId  string
	Processed bool
}

// Setup SQLite for local persistance
func InitSQLite() error {
	db, err := sql.Open("sqlite", "devpool.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Pool configuration
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	schema := `
CREATE TABLE IF NOT EXISTS read_issue_stream_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
  processed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS read_bounty_stream_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
  processed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS read_solution_stream_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
  processed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS read_achivement_stream_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
  processed BOOLEAN DEFAULT FALSE
);`

	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create tables: %w", err)
	}

	localDb = db
	return nil
}

func CloseDb() {
	if localDb != nil {
		if err := localDb.Close(); err != nil {
			Log.Error("Error closing database connection: %v", err)
		} else {
			Log.Info("Database connection closed.")
		}
	}
}

func ReadLastEntry(table string) (*Entry, error) {
	var tableName string
	switch table {
	case "bounty":
		tableName = "read_bounty_stream_events"
	case "issue":
		tableName = "read_issue_stream_events"
	case "solution":
		tableName = "read_solution_stream_events"
	case "achivement":
		tableName = "read_achivement_stream_events"
	default:
		return nil, fmt.Errorf("Invalid table name provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, stream_id, processed FROM ` + tableName + ` ORDER BY id DESC LIMIT 1`
	row := localDb.QueryRowContext(ctx, query)

	var entry Entry
	err := row.Scan(&entry.ID, &entry.StreamId, &entry.Processed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	Log.Info(
		fmt.Sprintf("Found Entry: ID=%d, StreamID=%s, Processed=%t\n",
			entry.ID,
			entry.StreamId,
			entry.Processed,
		),
	)
	return &entry, nil
}

func CreateEntry(table string, streamId string) (*Entry, error) {
	var tableName string
	switch table {
	case "bounty":
		tableName = "read_bounty_stream_events"
	case "issue":
		tableName = "read_issue_stream_events"
	case "solution":
		tableName = "read_solution_stream_events"
	case "achivement":
		tableName = "read_achivement_stream_events"
	default:
		return nil, fmt.Errorf("Invalid table name provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO ` + tableName + ` (stream_id) VALUES (?) RETURNING id, stream_id, processed`
	row := localDb.QueryRowContext(ctx, query, streamId)

	var entry Entry
	err := row.Scan(&entry.ID, &entry.StreamId, &entry.Processed)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}
	Log.Info(
		fmt.Sprintf("Created Entry: ID=%d, StreamID=%s, Processed=%t\n",
			entry.ID,
			entry.StreamId,
			entry.Processed,
		),
	)
	return &entry, nil
}

func UpdateEntryAsCompleted(table string, streamId string) (*Entry, error) {
	var tableName string
	switch table {
	case "bounty":
		tableName = "read_bounty_stream_events"
	case "issue":
		tableName = "read_issue_stream_events"
	case "solution":
		tableName = "read_solution_stream_events"
	case "achivement":
		tableName = "read_achivement_stream_events"
	default:
		return nil, fmt.Errorf("Invalid table name provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE ` + tableName + ` SET processed = TRUE WHERE stream_id = ? RETURNING id, stream_id, processed`
	row := localDb.QueryRowContext(ctx, query, streamId)

	var entry Entry
	err := row.Scan(&entry.ID, &entry.StreamId, &entry.Processed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("entry with StreamID %s not found", streamId)
		}
		return nil, err
	}

	Log.Info(fmt.Sprintf("Updated entry: ID=%d, StreamID=%s, Processed=%t\n", entry.ID, entry.StreamId, entry.Processed))
	return &entry, nil
}
