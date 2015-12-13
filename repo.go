package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
	"log"
	"os"
	"strings"
)

/*
* a repo[sitory] is the data access layer for the web application
 */

// dao is a collection of *sql.DB identified by strings (usernames), the point
// being that each user gets their own connection pool
var dao map[string]*sql.DB
var ddl goyesql.Queries
var dml goyesql.Queries
var conn string

// AppDB is an identifier for a specific *sql.DB in our dao map
const AppDB = "_app"

func init() {
	conn = os.Getenv("DATABASE_URL")

	if conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	dao = make(map[string]*sql.DB)
	ensureDb(AppDB)

	// Create all the Tables, Views if they do not exist
	ddl = goyesql.MustParseFile("data/sql/ddl.sql")

	logExec(dao[AppDB], (string)(ddl["create-user-table"]))
	logExec(dao[AppDB], (string)(ddl["create-command-table"]))
	logExec(dao[AppDB], (string)(ddl["create-tag-table"]))
	logExec(dao[AppDB], (string)(ddl["create-invocation-table"]))
	logExec(dao[AppDB], (string)(ddl["create-invocationtag-table"]))
	logExec(dao[AppDB], (string)(ddl["create-commandhistory-view"]))
	logExec(dao[AppDB], (string)(ddl["create-timestamp-index"]))

	// Load all data-manipulation queries
	dml = goyesql.MustParseFile("data/sql/queries.sql")

	log.Println("storage init completed")
}

func ql(identifier string) string {
	return (string)(dml[(goyesql.Tag)(identifier)])
}

func logExec(conn *sql.DB, query string) {
	log.Printf("executing query: '%s'\n", query)
	_, err := conn.Exec(query)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("executed query: '%s'\n", query)
	}
}

// ensureDb verifies that we have created a connection pool for the given user
func ensureDb(user string) {
	_, exists := dao[user]

	if !exists {
		var err error
		dao[user], err = sql.Open("postgres", conn)
		dao[user].SetMaxOpenConns(5)
		dao[user].SetMaxIdleConns(1)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// queryCommands is a common handler for implementing paging over Commands
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

// queryInvocations is a common handler for implementing paging over Invocations
func queryInvocations(rows *sql.Rows, pageSize int) (result Invocations, err error) {
	defer rows.Close()
	var tmp Invocation
	var tags string
	inc := 0
	if pageSize == 0 {
		inc = -1
	}
	for rows.Next() && inc < pageSize {
		err = rows.Scan(&tmp.Id, &tmp.ExitCode, &tmp.Timestamp,
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

// InsertInvocations sets up a transaction for commiting a batch of Invocations
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

// invocationTx handles the insertion of a single Invocation, as part of a batch
// transaction
func invocationTx(tx *sql.Tx, user string, inv Invocation) (err error) {
	var cmdid int
	err = tx.QueryRow(ql("get-commandid-by-command"), inv.Command).Scan(&cmdid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(ql("insert-command"), inv.Command).Scan(&cmdid)
		if err2 != nil {
			return err2
		}
	}

	var invid int
	err = tx.QueryRow(ql("insert-invocation"), user, cmdid, inv.ExitCode,
		inv.Timestamp, inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&invid)

	for _, tag := range inv.Tags {
		err = AddTag(tx, user, invid, tag)
	}
	return
}

// AddTag persists a Tag to an Invocation, as part of a transaction.
func AddTag(tx *sql.Tx, user string, invid int, tag string) (err error) {
	var tagid int
	err = tx.QueryRow(ql("get-tagid-by-name"), tag).Scan(&tagid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(ql("insert-tag"), tag).Scan(&tagid)
		if err2 != nil {
			return err2
		}
	}

	_, err = tx.Exec(ql("insert-invocationtag"), tagid, invid)
	return
}

// GetInvocations returns the [pageSize] most recent Invocations for the given
// user
func GetInvocations(user string, pageSize int) (result Invocations, err error) {
	ensureDb(user)
	rows, err := dao[user].Query(ql("get-invocations-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}

	return queryInvocations(rows, pageSize)
}

// GetCommands returns the [pageSize] most recent Commands for the given user
func GetCommands(user string, pageSize int) (result Commands, err error) {
	ensureDb(user)
	rows, err := dao[user].Query(ql("get-commands-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}
	return queryCommands(rows, pageSize)
}
