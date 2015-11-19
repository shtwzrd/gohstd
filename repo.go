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

	procs := []string{
		SqlInsertNotificationFunction,
		SqlCommandNotificationTrigger}

	for _, t := range tables {
		_, err := dao[AppDB].Exec(t)
		if err != nil {
			fmt.Println(t)
			log.Fatal(err)
		}
	}

	for _, t := range procs {
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
		dao[user].SetMaxOpenConns(90)
		dao[user].SetMaxIdleConns(0)

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

func InsertInvocations(user string, invocs Invocations) (err error) {
	ensureDb(user)
	tx, err := dao[user].Begin()
	if err != nil {
		return
	}
	for _, inv := range invocs {
		err = invocationTx(tx, user, inv)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func invocationTx(tx *sql.Tx, user string, inv Invocation) (err error) {
	var uid int
	err = tx.QueryRow(`
        SELECT userid FROM "user" WHERE username=$1`, user).Scan(&uid)

	var cmdid int
	err = tx.QueryRow(`
        SELECT commandid FROM command WHERE commandstring=$1`, inv.Command).
		Scan(&cmdid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(`
            INSERT INTO command (commandstring) VALUES ($1) RETURNING commandid`,
			inv.Command).Scan(&cmdid)
		if err2 != nil {
			return err2
		}
	}

	var ctxid int
	err = tx.QueryRow(`
        SELECT contextid FROM context
        WHERE hostname = $1 AND username = $2 AND shell = $3 AND directory = $4`,
		inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&ctxid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(`
            INSERT INTO context (hostname, username, shell, directory)
            VALUES ($1, $2, $3, $4) RETURNING contextid`,
			inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&ctxid)
		if err2 != nil {
			return err2
		}
	}

	var sessionid int
	err = tx.QueryRow(`
        SELECT sessionid FROM "session" WHERE contextid=$1`, ctxid).
		Scan(&sessionid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(`
            INSERT INTO "session" (contextid,"timestamp")
            VALUES($1, $2) RETURNING sessionid`,
			ctxid, inv.Timestamp).Scan(&sessionid)
		if err2 != nil {
			return err2
		}
	}

	var invid int
	err = tx.QueryRow(`
           INSERT INTO invocation
           (userid, commandid, returnstatus, "timestamp", sessionid)
           VALUES ($1, $2, $3, $4, $5) RETURNING invocationid`,
		uid, cmdid, inv.Status, inv.Timestamp, sessionid).Scan(&invid)

	for _, tag := range inv.Tags {
		err = addTag(tx, user, invid, tag)
	}
	return
}

func addTag(tx *sql.Tx, user string, invid int, tag string) (err error) {
	var tagid int
	err = tx.QueryRow(`
        SELECT tagid FROM tag WHERE name=$1`, tag).Scan(&tagid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(`
            INSERT INTO tag (name) VALUES ($1) RETURNING tagid`,
			tag).Scan(&tagid)
		if err2 != nil {
			return err2
		}
	}

	_, err = tx.Exec(`
           INSERT INTO invocationtag (tagid, invocationid)
           VALUES ($1, $2)`, tagid, invid)
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
