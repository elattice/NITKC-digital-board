package timetable

type Departure struct {
	ID            int    `json:"id"`
	Kind          string `json:"kind"`
	RouteName     string `json:"routeName"`
	Destination   string `json:"destination"`
	DepartureTime string `json:"departureTime"`
	Platform      string `json:"platform"`
}

type DeparturesResponse struct {
	Departures []Departure `json:"departures"`
}

type UpcomingDeparturesResponse struct {
	Bus       []Departure `json:"bus"`
	Train     []Departure `json:"train"`
	UpdatedAt string      `json:"updatedAt"`
}
