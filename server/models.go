package server

// Defines the data structuresto be fetched representing artists, locations, dates, and concert relations
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Date struct {
	Dates []string `json:"dates"`
}

type Loc struct {
	Locations []string `json:"locations"`
}

type Relation struct {
	DatesLocation map[string][]string `json:"datesLocations"`
}

// Passes dynamic data to HTML templates for rendering web pages.
type TemplateData struct {
	Title     string
	Artist    Artist
	Data      []Artist
	Locations Loc
	Dates     Date
	Concerts  Relation
	Query     string
	Results   []Artist
	Message   string
	Status    int
}
