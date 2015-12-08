package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	router := NewRouter()
	port := os.Getenv("PORT") // Get Heroku assigned PORT
	cert := os.Getenv("GOHST_CERT")
	key := os.Getenv("GOHST_KEY")

	if port == "" {
		port = "8080" // if not on Heroku or not defined, use 8080
	}

	if cert != "" && key != "" {
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, router))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}

}
