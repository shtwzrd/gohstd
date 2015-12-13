package main

import (
	"github.com/warreq/gohstd"
	"log"
	"net/http"
	"os"
)

func main() {

	conn := os.Getenv("DATABASE_URL")

	if conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}
	repo := *NewPsqlRepo(conn)

	port := os.Getenv("PORT") // Get Heroku assigned PORT
	cert := os.Getenv("GOHST_CERT")
	key := os.Getenv("GOHST_KEY")

	if port == "" {
		port = "8080" // if not on Heroku or not defined, use 8080
	}

	router := NewRouter(&repo)

	if cert != "" && key != "" {
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, router))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}

}
