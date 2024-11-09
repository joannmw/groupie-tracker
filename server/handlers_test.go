package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestRenderTemplate(t *testing.T) {
	// Mock templates
	templates = map[string]*template.Template{
		"test.html": template.Must(template.New("layout.html").Parse("{{.Title}}")),
	}

	tests := []struct {
		name     string
		tmpl     string
		data     interface{}
		expected int
	}{
		{"Valid template", "test.html", TemplateData{Title: "Test"}, http.StatusOK},
		{"Invalid template", "nonexistent.html", nil, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			renderTemplate(w, tt.tmpl, tt.data)
			if w.Code != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, w.Code)
			}
		})
	}
}

func TestCheckMethodAndPath(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedMethod string
		expectedPath   string
		expected       bool
	}{
		{"Correct method and path", "GET", "/test", "GET", "/test", true},
		{"Incorrect method", "POST", "/test", "GET", "/test", false},
		{"Incorrect path", "GET", "/wrong", "GET", "/test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.path, nil)
			result := checkMethodAndPath(w, r, tt.expectedMethod, tt.expectedPath)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMainPage(t *testing.T) {
	// Get the current working directory (server directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Move up one directory to project root
	projectRoot := filepath.Dir(wd)
	templateDir := filepath.Join(projectRoot, "templates")

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Fatalf("Template directory not found at %s", templateDir)
	}

	// Setup test cases
	tests := []struct {
		name          string
		method        string
		path          string
		expectedCode  int
		expectedTitle string
	}{
		{
			name:          "Valid GET request",
			method:        http.MethodGet,
			path:          "/",
			expectedCode:  http.StatusOK,
			expectedTitle: "Groupie Trackers - Artists",
		},
		{
			name:         "Wrong method",
			method:       http.MethodPost,
			path:         "/",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Wrong path",
			method:       http.MethodGet,
			path:         "/wrong",
			expectedCode: http.StatusNotFound,
		},
	}

	// Create template with full path
	layoutPath := filepath.Join(templateDir, "layout.html")
	indexPath := filepath.Join(templateDir, "index.html")

	// Verify template files exist
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		t.Fatalf("Layout template not found at %s", layoutPath)
	}
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatalf("Index template not found at %s", indexPath)
	}

	tmpl, err := template.ParseFiles(layoutPath, indexPath)
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	templates = map[string]*template.Template{
		"index.html": tmpl,
	}

	// Initialize artists slice with some test data
	artists = []Artist{
		{
			ID:   1,
			Name: "Test Artist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			MainPage(w, req)

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("MainPage() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			// For successful requests, check the response body
			if tt.expectedCode == http.StatusOK {
				// Check if title is in the response
				if !strings.Contains(w.Body.String(), tt.expectedTitle) {
					t.Errorf("MainPage() response doesn't contain expected title %v", tt.expectedTitle)
				}

				// Check if test artist data is in the response
				if !strings.Contains(w.Body.String(), "Test Artist") {
					t.Errorf("MainPage() response doesn't contain test artist data")
				}
			}
		})
	}
}

func TestInfoAboutArtist(t *testing.T) {
	// Get the current working directory (server directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Move up one directory to project root
	projectRoot := filepath.Dir(wd)
	templateDir := filepath.Join(projectRoot, "templates")

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Fatalf("Template directory not found at %s", templateDir)
	}

	// Create template with full path
	layoutPath := filepath.Join(templateDir, "layout.html")
	detailsPath := filepath.Join(templateDir, "details.html")

	// Verify template files exist
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		t.Fatalf("Layout template not found at %s", layoutPath)
	}
	if _, err := os.Stat(detailsPath); os.IsNotExist(err) {
		t.Fatalf("Details template not found at %s", detailsPath)
	}

	tmpl, err := template.ParseFiles(layoutPath, detailsPath)
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	templates = map[string]*template.Template{
		"details.html": tmpl,
	}

	// Initialize artists slice with some test data
	artists = []Artist{
		{
			ID:           1,
			Name:         "Test Artist",
			Locations:    "https://groupietrackers.herokuapp.com/api/locations/1",
			ConcertDates: "https://groupietrackers.herokuapp.com/api/dates/1",
			Relations:    "https://groupietrackers.herokuapp.com/api/relation/1",
		},
	}

	// Setup test cases
	tests := []struct {
		name           string
		method         string
		path           string
		query          string
		expectedCode   int
		expectedTitle  string
		expectedArtist string
	}{
		{
			name:           "Valid GET request",
			method:         http.MethodGet,
			path:           "/artists/",
			query:          "?id=1",
			expectedCode:   http.StatusOK,
			expectedTitle:  "Artist Details",
			expectedArtist: "Test Artist",
		},
		{
			name:         "Invalid ID",
			method:       http.MethodGet,
			path:         "/artists/",
			query:        "?id=999",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing ID",
			method:       http.MethodGet,
			path:         "/artists/",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Wrong method",
			method:       http.MethodPost,
			path:         "/artists/",
			query:        "?id=1",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest(tt.method, tt.path+tt.query, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			InfoAboutArtist(w, req)

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("InfoAboutArtist() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			// For successful requests, check the response body
			if tt.expectedCode == http.StatusOK {
				// Check if title is in the response
				if !strings.Contains(w.Body.String(), tt.expectedTitle) {
					t.Errorf("InfoAboutArtist() response doesn't contain expected title %v", tt.expectedTitle)
				}

				// Check if test artist data is in the response
				if !strings.Contains(w.Body.String(), tt.expectedArtist) {
					t.Errorf("InfoAboutArtist() response doesn't contain expected artist %v", tt.expectedArtist)
				}
			}
		})
	}
}

func TestSearchPage(t *testing.T) {
	// Get the current working directory (server directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Move up one directory to project root
	projectRoot := filepath.Dir(wd)
	templateDir := filepath.Join(projectRoot, "templates")

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Fatalf("Template directory not found at %s", templateDir)
	}

	// Create template with full path
	layoutPath := filepath.Join(templateDir, "layout.html")
	searchPath := filepath.Join(templateDir, "search.html")

	// Verify template files exist
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		t.Fatalf("Layout template not found at %s", layoutPath)
	}
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		t.Fatalf("Search template not found at %s", searchPath)
	}

	tmpl, err := template.ParseFiles(layoutPath, searchPath)
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	templates = map[string]*template.Template{
		"search.html": tmpl,
	}

	// Initialize artists slice with some test data
	artists = []Artist{
		{
			ID:   1,
			Name: "Test Artist",
		},
		{
			ID:   2,
			Name: "Another Artist",
		},
	}

	// Setup test cases
	tests := []struct {
		name            string
		method          string
		path            string
		query           string
		expectedCode    int
		expectedTitle   string
		expectedQuery   string
		expectedArtist  string
		expectedMessage string
	}{
		{
			name:           "Valid search query",
			method:         http.MethodGet,
			path:           "/search/",
			query:          "?q=test",
			expectedCode:   http.StatusOK,
			expectedTitle:  "Search Results",
			expectedQuery:  "test",
			expectedArtist: "Test Artist",
		},
		{
			name:            "No results found",
			method:          http.MethodGet,
			path:            "/search/",
			query:           "?q=nonexistent",
			expectedCode:    http.StatusOK,
			expectedTitle:   "Search Results",
			expectedQuery:   "nonexistent",
			expectedMessage: "No matching artists found.",
		},
		{
			name:         "Empty query",
			method:       http.MethodGet,
			path:         "/search/",
			query:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Wrong method",
			method:       http.MethodPost,
			path:         "/search/",
			query:        "?q=test",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest(tt.method, tt.path+tt.query, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			SearchPage(w, req)

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("SearchPage() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			// For successful requests, check the response body
			if tt.expectedCode == http.StatusOK {
				body := w.Body.String()

				// Check if title is in the response
				if !strings.Contains(body, tt.expectedTitle) {
					t.Errorf("SearchPage() response doesn't contain expected title %v", tt.expectedTitle)
				}

				// Check if query is in the response
				if tt.expectedQuery != "" && !strings.Contains(body, tt.expectedQuery) {
					t.Errorf("SearchPage() response doesn't contain expected query %v", tt.expectedQuery)
				}

				// Check if test artist data is in the response for a valid query
				if tt.expectedArtist != "" && !strings.Contains(body, tt.expectedArtist) {
					t.Errorf("SearchPage() response doesn't contain expected artist %v", tt.expectedArtist)
				}

				// Check for "no results" message if expected
				if tt.expectedMessage != "" && !strings.Contains(body, tt.expectedMessage) {
					t.Errorf("SearchPage() response doesn't contain expected message %v", tt.expectedMessage)
				}
			}
		})
	}
}

func TestErrorPage(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{"Not Found", http.StatusNotFound, "404 - Not Found"},
		{"Bad Request", http.StatusBadRequest, "400 - Bad Request"},
		{"Method Not Allowed", http.StatusMethodNotAllowed, "405 - Method Not Allowed"},
		{"Internal Server Error", http.StatusInternalServerError, "500 - Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ErrorPage(w, tt.code)
			if w.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}
			if !strings.HasPrefix(w.Body.String(), tt.expected) {
				t.Errorf("Expected body %s, got %s,", tt.expected, w.Body.String())
			}
		})
	}
}

func TestServeStatic_Success(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Could not change directory: %v", err)
	}
	// Ensure we change back to the original directory after the test
	defer func() {
		err := os.Chdir("server")
		if err != nil {
			t.Fatalf("Could not change back to original directory: %v", err)
		}
	}()
	// Create a response recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/static/style.css", nil)
	// Call the handler function
	ServeStatic(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	responseBody := w.Body.Bytes()
	// Read the expected content from the file
	expectedContent, err := os.ReadFile("static/style.css")
	if err != nil {
		t.Fatalf("Failed to read expected content from file: %v", err)
	}
	if !bytes.Equal(responseBody, expectedContent) {
		t.Errorf("Expected response body to match the image file content")
	}
}

func TestServeStatic_Forbidden(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Could not change directory: %v", err)
	}
	// Ensure we change back to the original directory after the test
	defer func() {
		err := os.Chdir("server")
		if err != nil {
			t.Fatalf("Could not change back to original directory: %v", err)
		}
	}()
	// Create a response recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/static/nonexistent.txt", nil)
	ServeStatic(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
	responseBody := w.Body.String()
	expectedErrorMessage := "Not Found"
	if !strings.Contains(responseBody, expectedErrorMessage) {
		t.Errorf("Expected response body to contain '%s', but it didn't", expectedErrorMessage)
	}
}

func TestServeStatic_DirectoryHandling(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Could not change directory: %v", err)
	}
	// Ensure we change back to the original directory after the test
	defer func() {
		err := os.Chdir("server")
		if err != nil {
			t.Fatalf("Could not change back to original directory: %v", err)
		}
	}()
	// Create a response recorder
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/static/", nil)
	ServeStatic(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
	responseBody := w.Body.String()
	expectedErrorMessage := "Not Found"
	if !strings.Contains(responseBody, expectedErrorMessage) {
		t.Errorf("Expected response body to contain '%s', but it didn't", expectedErrorMessage)
	}
}
