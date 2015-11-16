package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var dao map[string]*sql.DB
var conn string

const AppDB = "_app"

func init() {
	conn = os.Getenv("DATABASE_URL")

	if conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	dao = make(map[string]*sql.DB)
	ensureDb(AppDB)

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
		_, err := dao[AppDB].Exec(t)
		if err != nil {
			fmt.Println(t)
			log.Fatal(err)
		}
	}
}

func ensureDb(user string) {
	var err error
	dao[user], err = sql.Open("postgres", conn)

	if err != nil {
		log.Fatal(err)
	}
}

func queryCommands(rows *sql.Rows, pageSize int) (result Commands, err error) {
	defer rows.Close()
	var c string
	var inc int
	for rows.Next() && (inc < pageSize || pageSize == 0) {
		err = rows.Scan(&c)
		if err != nil {
			return
		}
		result = append(result, Command(c))
		inc++
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

func queryInvocations(rows *sql.Rows, pageSize int) (result Invocations, err error) {
	defer rows.Close()
	var tmp Invocation
	var tags string
	var inc int
	var x interface{} // for ignoring scan columns
	for rows.Next() && (inc < pageSize || pageSize == 0) {
		err = rows.Scan(&x, &tmp.Id, &tmp.SessionId, &tmp.Status, &tmp.Timestamp,
			&tmp.Host, &tmp.User, &tmp.Shell, &tmp.Directory, &tmp.Command, &tags)
		if err != nil {
			log.Println(err)
			return
		}
		tmp.Tags = strings.Split(tags[1:len(tags)-1], ", ")
		result = append(result, tmp)
		inc++
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

func GetAllInvocations(user string) (result Invocations, err error) {
	query := `select * from commandhistory where "user" = $1`

	ensureDb(user)
	rows, err := dao[user].Query(query, user)
	if err != nil {
		log.Println(err)
		return
	}

	return queryInvocations(rows, 0)
}

func GetAllCommands(user string) (result Commands, err error) {
	query := `select commandstring from commandhistory where "user" = $1`

	ensureDb(user)
	rows, err := dao[user].Query(query, user)
	if err != nil {
		log.Println(err)
		return
	}
	return queryCommands(rows, 0)
}
