package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// renderTemplate renders a specified template with the provided data.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Retrieve the template from the global map
	t, ok := templates[tmpl]
	if !ok {
		log.Println(tmpl, "not found")
		ErrorPage(w, http.StatusNotFound)
		return
	}
	// Execute the template with the provided data and layout
	err := t.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
}

// checkMethodAndPath checks if the request method and path match expected values.
func checkMethodAndPath(w http.ResponseWriter, r *http.Request, method, path string) bool {
	// Render a 405 error page for wrong method
	if r.Method != method {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return false
	}
	// Render a 404 error page for wrong path
	if r.URL.Path != path {
		ErrorPage(w, http.StatusNotFound)
		return false
	}
	return true
}

// MainPage serves as the home page of the application.
func MainPage(w http.ResponseWriter, r *http.Request) {
	if !checkMethodAndPath(w, r, http.MethodGet, "/") {
		return
	}
	// Create a TemplateData object with the title and list of artists.
	data := TemplateData{
		Title: "Groupie Trackers - Artists",
		Data:  artists,
	}
	renderTemplate(w, "index.html", data)
}

// InfoAboutArtist serves detailed information about a specific artist.
func InfoAboutArtist(w http.ResponseWriter, r *http.Request) {
	if !checkMethodAndPath(w, r, http.MethodGet, "/artists/") {
		return
	}

	// Get artist ID through query parameter and validate
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 || id > len(artists) || err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusBadRequest)
		return
	}
	id--

	// Fetch artist data
	locations, err := FetchLocations(artists[id].Locations)
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	dates, err := FetchDates(artists[id].ConcertDates)
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	rel, err := FetchRelation(artists[id].Relations)
	if err != nil {
		log.Println(err)
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Title:     "Artist Details",
		Artist:    artists[id],
		Locations: locations,
		Dates:     dates,
		Concerts:  rel,
	}
	// Render the artist details template with all relevant data
	renderTemplate(w, "details.html", data)
}

// SearchPage handles the artist search functionality.
func SearchPage(w http.ResponseWriter, r *http.Request) {
	if !checkMethodAndPath(w, r, http.MethodGet, "/search/") {
		return
	}

	// Get search query from URL parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	var results []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			results = append(results, artist)
		}
	}

	data := TemplateData{
		Title:   "Search Results",
		Query:   query,
		Results: results,
	}

	if len(results) == 0 {
		data.Message = "No artists found matching your query."
	}

	// Render the search results template with matched artists
	renderTemplate(w, "search.html", data)
}

// ErrorPage renders an error page based on the HTTP status code.
func ErrorPage(w http.ResponseWriter, code int) {
	var message string
	switch code {
	case http.StatusNotFound:
		message = "Not Found"
	case http.StatusBadRequest:
		message = "Bad Request"
	case http.StatusMethodNotAllowed:
		message = "Method Not Allowed"
	case http.StatusForbidden:
		message = "Forbidden"
	default:
		message = "Internal Server Error"
	}
	data := TemplateData{
		Title:   "Error",
		Status:  code,
		Message: message,
	}

	// Set HTTP response status code
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("templates/errors.html")
	// Serve basic error response if template parsing fails
	if err != nil {
		http.Error(w, fmt.Sprintf("%d - %s", code, message), code)
		return
	}

	// Serve basic error response if template execution fails
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%d - %s", code, message), code)
	}
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}
	// Remove the /static/ prefix from the URL path
	filePath := path.Join("static", strings.TrimPrefix(r.URL.Path, "/static/"))

	// Check if the file exists and is not a directory
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		ErrorPage(w, http.StatusNotFound)
		return
	}

	// Check the file extension
	ext := filepath.Ext(filePath)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".otf":
		w.Header().Set("Content-Type", "font/otf")
	default:
		ErrorPage(w, http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}
