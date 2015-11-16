package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dao map[string]*sql.DB

const AppDB = "_app"

func init() {
	conn := os.Getenv("DATABASE_URL")

	if conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	var err error
	dao = make(map[string]*sql.DB)
	dao[AppDB], err = sql.Open("postgres", conn)

	if err != nil {
		log.Fatal(err)
	}

	// Create all the Tables, Views if they do not exist
	tables := []string{
		SqlCreateUserTable,
		SqlCreateCommandTable,
		SqlCreateConfigurationTable,
		SqlCreateContextTable,
		SqlCreateSessionTable,
		SqlCreateInvocationTable,
		SqlCreateTagTable,
		SqlCreateInvocationTagTable,
		SqlCreateNotificationTable,
		SqlCreateServiceLogTable,
		SqlCreateCommandHistoryView}

	for _, t := range tables {
		_, err = dao[AppDB].Exec(t)
		if err != nil {
			fmt.Println(t)
			log.Fatal(err)
		}
	}
}
