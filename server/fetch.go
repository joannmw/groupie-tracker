package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// Global variables to hold templates and artist data
var templates map[string]*template.Template
var artists []Artist

var artistsURL = "https://groupietrackers.herokuapp.com/api/artists"

// init initializes templates and fetches artist data when the package is loaded.
func init() {
	var err error
	// Load HTML templates into the templates map
	templates, err = loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	if err := FetchArtists(); err != nil {
		log.Fatal("could not fetch artists: ", err)
	}
}

// loadTemplates loads HTML templates from the templates directory.
func loadTemplates() (map[string]*template.Template, error) {
	templates = make(map[string]*template.Template)
	layout := "templates/layout.html"

	// Get all HTML files in the "templates" directory
	pages, err := filepath.Glob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to load template files: %w", err)
	}

	for _, page := range pages {
		if page == layout {
			continue
		}
		// Combine layout with the current page template and parse the templates
		files := []string{layout, page}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", page, err)
		}
		//Store the parsed template in the map using the base file name as key
		templates[filepath.Base(page)] = tmpl
	}
	return templates, nil
}

// Fetches data from the given URL and unmarshals it into the target struct.
func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}

	defer response.Body.Close()

	// Read response body into bytes slice
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	// Unmarshal JSON into target struct
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal data from %s: %w", url, err)
	}

	return nil
}

// FetchArtists retrieves artist data from the defined artistsURL using fetchData
func FetchArtists() error {
	return fetchData(artistsURL, &artists)
}

// FetchRelation retrieves relation data from a specified URL using fetchData
func FetchLocations(url string) (Loc, error) {
	var location Loc
	err := fetchData(url, &location)
	return location, err
}

// FetchRelation retrieves relation data from a specified URL using fetchData
func FetchRelation(url string) (Relation, error) {
	var relation Relation
	err := fetchData(url, &relation)
	return relation, err
}

// FetchDates retrieves date data from a specified URL using fetchData
func FetchDates(url string) (Date, error) {
	var dates Date
	err := fetchData(url, &dates)
	return dates, err
}
