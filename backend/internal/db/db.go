package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS departures (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  kind TEXT NOT NULL CHECK (kind IN ('bus', 'train')),
  route_name TEXT NOT NULL CHECK (route_name <> ''),
  destination TEXT NOT NULL CHECK (destination <> ''),
  departure_time TEXT NOT NULL CHECK (departure_time <> ''),
  platform TEXT NOT NULL DEFAULT '',
  active INTEGER NOT NULL DEFAULT 1
);
`

const migratedSchema = `
CREATE TABLE departures_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  kind TEXT NOT NULL CHECK (kind IN ('bus', 'train')),
  route_name TEXT NOT NULL CHECK (route_name <> ''),
  destination TEXT NOT NULL CHECK (destination <> ''),
  departure_time TEXT NOT NULL CHECK (departure_time <> ''),
  platform TEXT NOT NULL DEFAULT '',
  active INTEGER NOT NULL DEFAULT 1
);
`

// Open opens the SQLite database at path, creating the parent directory and
// schema as needed. Timetable data is only populated through CSV import.
func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// modernc.org/sqlite does not support concurrent writers on one file;
	// a single connection avoids SQLITE_BUSY errors.
	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := migrateSchema(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return conn, nil
}

func migrateSchema(conn *sql.DB) error {
	rows, err := conn.Query(`PRAGMA table_info(departures)`)
	if err != nil {
		return err
	}

	columns := map[string]bool{}
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue any
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			rows.Close()
			return err
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	if len(columns) == 0 {
		_, err := conn.Exec(schema)
		return err
	}
	if !columns["origin"] && !columns["arrival_time"] && !columns["note"] {
		return nil
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migratedSchema); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		INSERT INTO departures_new (id, kind, route_name, destination, departure_time, platform, active)
		SELECT id, kind, route_name, destination, departure_time, COALESCE(platform, ''), active
		FROM departures
	`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE departures`); err != nil {
		return err
	}
	if _, err := tx.Exec(`ALTER TABLE departures_new RENAME TO departures`); err != nil {
		return err
	}
	return tx.Commit()
}
