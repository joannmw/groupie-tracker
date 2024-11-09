package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	// Create temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "test_templates")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create "templates" subdirectory
	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.Mkdir(templatesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create test template files
	layoutPath := filepath.Join(templatesDir, "layout.html")
	pagePath := filepath.Join(templatesDir, "page.html")

	if err := os.WriteFile(layoutPath, []byte("layout"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pagePath, []byte("page"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set working directory to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Test loadTemplates
	templates, err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates() error = %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("loadTemplates() returned %d templates, want 1", len(templates))
	}
	if _, ok := templates["page.html"]; !ok {
		t.Errorf("loadTemplates() did not load page.html")
	}
}

func TestFetchData(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name": "Test Artist"}`)
	}))
	defer ts.Close()

	// Test fetchData
	var result map[string]string
	err := fetchData(ts.URL, &result)
	if err != nil {
		t.Fatalf("fetchData() error = %v", err)
	}
	expected := map[string]string{"name": "Test Artist"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("fetchData() = %v, want %v", result, expected)
	}
}

func TestFetchArtists(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Artist{{Name: "Test Artist"}})
	}))
	defer ts.Close()

	// Replace the real URL with the test server URL
	oldURL := artistsURL
	artistsURL = ts.URL
	defer func() { artistsURL = oldURL }()

	// Test FetchArtists
	err := FetchArtists()
	if err != nil {
		t.Fatalf("FetchArtists() error = %v", err)
	}
	if len(artists) != 1 || artists[0].Name != "Test Artist" {
		t.Errorf("FetchArtists() did not populate artists correctly")
	}
}

func TestFetchLocations(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Loc{Locations: []string{"Test Location"}})
	}))
	defer ts.Close()

	// Test FetchLocations
	loc, err := FetchLocations(ts.URL)
	if err != nil {
		t.Fatalf("FetchLocations() error = %v", err)
	}
	if len(loc.Locations) != 1 || loc.Locations[0] != "Test Location" {
		t.Errorf("FetchLocations() returned incorrect data")
	}
}

func TestFetchRelation(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Relation{DatesLocation: map[string][]string{"Test": {"2023"}}})
	}))
	defer ts.Close()

	// Test FetchRelation
	relation, err := FetchRelation(ts.URL)
	if err != nil {
		t.Fatalf("FetchRelation() error = %v", err)
	}
	if len(relation.DatesLocation) != 1 || len(relation.DatesLocation["Test"]) != 1 || relation.DatesLocation["Test"][0] != "2023" {
		t.Errorf("FetchRelation() returned incorrect data")
	}
}

func TestFetchDates(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Date{Dates: []string{"2023-01-01"}})
	}))
	defer ts.Close()

	// Test FetchDates
	dates, err := FetchDates(ts.URL)
	if err != nil {
		t.Fatalf("FetchDates() error = %v", err)
	}
	if len(dates.Dates) != 1 || dates.Dates[0] != "2023-01-01" {
		t.Errorf("FetchDates() returned incorrect data")
	}
}
