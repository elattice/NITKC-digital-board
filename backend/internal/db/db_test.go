package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOpenLeavesNewDatabaseEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timetable.db")

	conn, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopening simulates a server restart before any CSV is imported.
	conn, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM departures`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("departure count = %d, want 0", count)
	}
}

func TestOpenMigratesLegacyDeparturesSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timetable.db")
	legacy, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = legacy.Exec(`
		CREATE TABLE departures (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kind TEXT NOT NULL,
			route_name TEXT NOT NULL,
			origin TEXT NOT NULL,
			destination TEXT NOT NULL,
			departure_time TEXT NOT NULL,
			arrival_time TEXT,
			platform TEXT,
			note TEXT,
			active INTEGER NOT NULL DEFAULT 1
		);
		INSERT INTO departures
			(kind, route_name, origin, destination, departure_time, arrival_time, platform, note)
		VALUES
			('bus', 'テスト路線', '高専前', '大楽毛駅', '10:15', '10:35', '高専前', '備考');
	`)
	if err != nil {
		legacy.Close()
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}

	conn, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	rows, err := conn.Query(`PRAGMA table_info(departures)`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue any
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		columns = append(columns, name)
	}
	wantColumns := []string{"id", "kind", "route_name", "destination", "departure_time", "platform", "active"}
	if !reflect.DeepEqual(columns, wantColumns) {
		t.Fatalf("columns = %v, want %v", columns, wantColumns)
	}

	var kind, routeName, destination, departureTime, platform string
	err = conn.QueryRow(`
		SELECT kind, route_name, destination, departure_time, platform
		FROM departures
	`).Scan(&kind, &routeName, &destination, &departureTime, &platform)
	if err != nil {
		t.Fatal(err)
	}
	if kind != "bus" || routeName != "テスト路線" || destination != "大楽毛駅" || departureTime != "10:15" || platform != "高専前" {
		t.Fatalf("migrated row = %q, %q, %q, %q, %q", kind, routeName, destination, departureTime, platform)
	}
}
