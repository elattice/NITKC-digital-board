package timetable

import (
	"database/sql"
	"time"
)

// Service returns timetable data stored in SQLite.
type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// UpcomingDepartures returns active departures at or after the current time,
// grouped by kind for the board. departure_time is a zero-padded HH:MM string,
// so lexicographic comparison matches chronological order within the day.
func (s *Service) UpcomingDepartures() (UpcomingDeparturesResponse, error) {
	now := time.Now().Format("15:04")
	departures, err := s.queryDepartures(`
		SELECT id, kind, route_name, destination, departure_time, platform
		FROM departures
		WHERE active = 1 AND departure_time >= ?
		ORDER BY departure_time, id
	`, now)
	if err != nil {
		return UpcomingDeparturesResponse{}, err
	}

	response := UpcomingDeparturesResponse{
		Bus:       []Departure{},
		Train:     []Departure{},
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	for _, d := range departures {
		switch d.Kind {
		case "bus":
			response.Bus = append(response.Bus, d)
		case "train":
			response.Train = append(response.Train, d)
		}
	}
	return response, nil
}

// AllDepartures returns every departure, including inactive ones.
func (s *Service) AllDepartures() ([]Departure, error) {
	return s.queryDepartures(`
		SELECT id, kind, route_name, destination, departure_time, platform
		FROM departures
		ORDER BY departure_time, id
	`)
}

// ReplaceDepartures deletes all departures and inserts the given ones,
// in a single transaction so a failure leaves the current data intact.
func (s *Service) ReplaceDepartures(departures []Departure) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM departures`); err != nil {
		return 0, err
	}
	// Restart AUTOINCREMENT ids from 1 for the replaced data.
	if _, err := tx.Exec(`DELETE FROM sqlite_sequence WHERE name = 'departures'`); err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO departures (kind, route_name, destination, departure_time, platform)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, d := range departures {
		if _, err := stmt.Exec(d.Kind, d.RouteName, d.Destination, d.DepartureTime, d.Platform); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(departures), nil
}

func (s *Service) queryDepartures(query string, args ...any) ([]Departure, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departures := []Departure{}
	for rows.Next() {
		var d Departure
		if err := rows.Scan(&d.ID, &d.Kind, &d.RouteName, &d.Destination, &d.DepartureTime, &d.Platform); err != nil {
			return nil, err
		}
		departures = append(departures, d)
	}
	return departures, rows.Err()
}
