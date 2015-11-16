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
	inc := 0
	if pageSize == 0 {
		inc = -1
	}
	for rows.Next() && inc < pageSize {
		err = rows.Scan(&c)
		if err != nil {
			return
		}
		result = append(result, Command(c))
		if pageSize > 0 {
			inc++
		}
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
	inc := 0
	if pageSize == 0 {
		inc = -1
	}
	for rows.Next() && inc < pageSize {
		err = rows.Scan(&tmp.Id, &tmp.SessionId, &tmp.Status, &tmp.Timestamp,
			&tmp.Host, &tmp.User, &tmp.Shell, &tmp.Directory, &tmp.Command, &tags)
		if err != nil {
			log.Println(err)
			return
		}
		tmp.Tags = strings.Split(tags[1:len(tags)-1], ", ")
		result = append(result, tmp)
		if pageSize > 0 {
			inc++
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

func GetInvocations(user string, pageSize int) (result Invocations, err error) {
	query := `SELECT invocationid, sessionid, returnstatus, "timestamp", hostname,
            username, shell, directory, commandstring, tags
            FROM commandhistory WHERE "user" = $1`

	ensureDb(user)
	rows, err := dao[user].Query(query, user)
	if err != nil {
		log.Println(err)
		return
	}

	return queryInvocations(rows, pageSize)
}

func GetCommands(user string, pageSize int) (result Commands, err error) {
	query := `SELECT commandstring FROM commandhistory WHERE "user" = $1`

	ensureDb(user)
	rows, err := dao[user].Query(query, user)
	if err != nil {
		log.Println(err)
		return
	}
	return queryCommands(rows, pageSize)
}
