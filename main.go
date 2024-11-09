package main

import (
	"fmt"
	"log"
	"net/http"

	"groupie-tracker/server"
)

func main() {
	http.HandleFunc("/static/", server.ServeStatic)
	http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/artists/", server.InfoAboutArtist)
	http.HandleFunc("/search/", server.SearchPage)
	fmt.Println("Server running on http://localhost:3000/")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
