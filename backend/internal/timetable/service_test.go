package timetable

import (
	"path/filepath"
	"testing"

	"open-campus-board/backend/internal/db"
)

func TestReplaceDeparturesRollsBackOnInsertError(t *testing.T) {
	conn, err := db.Open(filepath.Join(t.TempDir(), "timetable.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	service := NewService(conn)
	baseline := []Departure{{
		Kind:          "bus",
		RouteName:     "既存路線",
		Destination:   "既存行き先",
		DepartureTime: "10:00",
		Platform:      "1番",
	}}
	if _, err := service.ReplaceDepartures(baseline); err != nil {
		t.Fatal(err)
	}

	invalid := []Departure{
		{
			Kind:          "train",
			RouteName:     "新しい路線",
			Destination:   "新しい行き先",
			DepartureTime: "11:00",
			Platform:      "",
		},
		{
			Kind:          "invalid",
			RouteName:     "不正路線",
			Destination:   "不正行き先",
			DepartureTime: "12:00",
			Platform:      "",
		},
	}
	if _, err := service.ReplaceDepartures(invalid); err == nil {
		t.Fatal("ReplaceDepartures() error = nil, want insert error")
	}

	departures, err := service.AllDepartures()
	if err != nil {
		t.Fatal(err)
	}
	if len(departures) != 1 || departures[0].RouteName != "既存路線" {
		t.Fatalf("departures after rollback = %#v, want original row", departures)
	}
}
