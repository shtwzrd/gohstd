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
	_, exists := dao[user]

	if !exists {
		var err error
		dao[user], err = sql.Open("postgres", conn)

		if err != nil {
			log.Fatal(err)
		}
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

func InsertInvocation(user string, inv Invocation) (err error) {
	ensureDb(user)
	tx, err := dao[user].Begin()
	if err != nil {
		return
	}

	var uid int
	err = tx.QueryRow(`SELECT userid FROM "user" WHERE username=$1`,
		user).Scan(&uid)
	if err != nil {
		return
	}

	var cmdid int
	err = tx.QueryRow(`SELECT commandid FROM command WHERE commandstring=$1`,
		inv.Command).Scan(&cmdid)

	switch {
	case err == sql.ErrNoRows:
		err2 := tx.QueryRow(`INSERT INTO command (commandstring)
 VALUES ($1) RETURNING commandid`, inv.Command).Scan(&cmdid)
		if err2 != nil {
			return err2
		}
		break
	case err != nil:
		tx.Rollback()
		return
	}

	var ctxid int
	err = tx.QueryRow(`SELECT contextid FROM context WHERE
 hostname = $1 AND username = $2 AND shell = $3 AND directory = $4`,
		inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&ctxid)

	switch {
	case err == sql.ErrNoRows:
		err2 := tx.QueryRow(`INSERT INTO context
 (hostname, username, shell, directory) VALUES ($1, $2, $3, $4)
 RETURNING contextid`,
			inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&ctxid)
		if err2 != nil {
			return err2
		}
		break
	case err != nil:
		tx.Rollback()
		return
	}

	var sessionid int
	err = tx.QueryRow(`SELECT sessionid FROM "session" WHERE
 contextid=$1`, ctxid).Scan(&sessionid)

	switch {
	case err == sql.ErrNoRows:
		err2 := tx.QueryRow(`INSERT INTO "session" (contextid,"timestamp")
 VALUES($1, $2) RETURNING sessionid`, ctxid, inv.Timestamp).Scan(&sessionid)
		if err2 != nil {
			return err2
		}
		break
	case err != nil:
		tx.Rollback()
		return
	}

	// Finally insert the invocation
	_, err = tx.Exec(`INSERT INTO invocation
 (userid, commandid, returnstatus, "timestamp", sessionid)
 VALUES ($1, $2, $3, $4, $5)`, uid, cmdid, inv.Status, inv.Timestamp, sessionid)

	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
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
