package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	router := NewRouter()
	port := os.Getenv("PORT") // Get Heroku assigned PORT

	if port == "" {
		port = "8080" // if not on Heroku or not defined, use 8080
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
