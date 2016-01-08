package service

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
	"log"
	"os"
)

// PsqlDao is an entity that implements per-user connection pool limits and
// facilitates access to SQL queries loaded from disk.
type PsqlDao struct {
	// dao is a collection of *sql.DB identified by strings (usernames), the point
	// being that each user gets their own connection pool
	dao  map[string]*sql.DB
	ddl  goyesql.Queries
	dml  goyesql.Queries
	conn string
}

// AppDB is an identifier for a specific *sql.DB in our dao map
const AppDB = "_app"

func NewPsqlDao(conn string) *PsqlDao {
	db := PsqlDao{}
	db.conn = os.Getenv("DATABASE_URL")

	if db.conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	db.dao = make(map[string]*sql.DB)

	// Create all the Tables, Views if they do not exist
	db.ddl = goyesql.MustParseFile("data/sql/ddl.sql")

	db.EnsurePool(AppDB)

	logExec(db.dao[AppDB], (string)(db.ddl["create-user-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-post-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-command-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-tag-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-invocation-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-invocationtag-table"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-commandhistory-view"]))
	logExec(db.dao[AppDB], (string)(db.ddl["create-timestamp-index"]))

	// Load all data-manipulation queries
	db.dml = goyesql.MustParseFile("data/sql/queries.sql")

	log.Println("storage init completed")
	return &db
}

// Query fetches an SQL query that matches the given identifier
func (r *PsqlDao) Query(identifier string) string {
	return (string)(r.dml[(goyesql.Tag)(identifier)])
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

// EnsurePool verifies that we have created a connection pool for the given user
func (r *PsqlDao) EnsurePool(user string) *sql.DB {
	_, exists := r.dao[user]

	if !exists {
		var err error
		r.dao[user], err = sql.Open("postgres", r.conn)
		r.dao[user].SetMaxOpenConns(2)
		r.dao[user].SetMaxIdleConns(0)

		if err != nil {
			log.Fatal(err)
		}
	}
	return r.dao[user]
}
