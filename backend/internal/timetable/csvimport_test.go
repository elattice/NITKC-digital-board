package timetable

import (
	"strings"
	"testing"
)

func TestParseCSVAcceptsSimplifiedFormat(t *testing.T) {
	input := `kind,route_name,destination,departure_time,platform
bus,高専前→大楽毛駅,大楽毛駅,09:40,高専前
train,大楽毛駅→釧路方面,釧路方面,10:22,
`

	departures, errs := ParseCSV(strings.NewReader(input))
	if len(errs) != 0 {
		t.Fatalf("ParseCSV() errors = %v", errs)
	}
	if len(departures) != 2 {
		t.Fatalf("len(departures) = %d, want 2", len(departures))
	}
	if departures[0].Platform != "高専前" {
		t.Errorf("bus platform = %q, want %q", departures[0].Platform, "高専前")
	}
	if departures[1].Platform != "" {
		t.Errorf("train platform = %q, want empty", departures[1].Platform)
	}
}

func TestParseCSVValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		row     string
		message string
	}{
		{name: "kind", row: "airplane,路線,行き先,10:00,", message: "kind は bus または train"},
		{name: "route name", row: "train,,行き先,10:00,", message: "route_name が空です"},
		{name: "destination", row: "train,路線,,10:00,", message: "destination が空です"},
		{name: "departure time", row: "train,路線,行き先,9:00,", message: "departure_time は HH:MM 形式"},
		{name: "bus platform", row: "bus,路線,行き先,10:00,", message: "bus の platform が空です"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "kind,route_name,destination,departure_time,platform\n" + tt.row + "\n"
			_, errs := ParseCSV(strings.NewReader(input))
			if !containsMessage(errs, tt.message) {
				t.Fatalf("ParseCSV() errors = %v, want message containing %q", errs, tt.message)
			}
		})
	}
}

func TestParseCSVRejectsOldHeader(t *testing.T) {
	input := "kind,route_name,origin,destination,departure_time,arrival_time,platform,note\n"
	_, errs := ParseCSV(strings.NewReader(input))
	if !containsMessage(errs, "ヘッダー行が不正です") {
		t.Fatalf("ParseCSV() errors = %v, want invalid header error", errs)
	}
}

func containsMessage(messages []string, want string) bool {
	for _, message := range messages {
		if strings.Contains(message, want) {
			return true
		}
	}
	return false
}
