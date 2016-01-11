package main

import (
	gohst "github.com/warreq/gohstd/service"
	"log"
	"net/http"
	"os"
)

func main() {

	conn := os.Getenv("DATABASE_URL")

	if conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}
	dao := *gohst.NewPsqlDao(conn)
	cmd := *gohst.NewPsqlCommandRepo(dao)
	user := *gohst.NewPsqlUserRepo(dao)
	post := *gohst.NewPsqlPostRepo(dao)

	port := os.Getenv("PORT") // Get Heroku assigned PORT
	httpOnlyPort := os.Getenv("HTTP_ONLY_PORT")
	cert := os.Getenv("GOHST_CERT")
	key := os.Getenv("GOHST_KEY")

	if port == "" {
		port = "8080" // if not on Heroku or not defined, use 8080
	}

	router := gohst.NewRouter(cmd, user, post)

	if cert != "" && key != "" {
		if httpOnlyPort != "" {
			go log.Fatal(http.ListenAndServe(":"+httpOnlyPort, nil))
		}
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, router))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}

}
