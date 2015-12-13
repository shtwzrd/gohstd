package gohstd

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
	"log"
	"os"
	"strings"
)

// PsqlRepo is an implementation of a gohst Repo that uses PostgreSQL as a
// backing store
type PsqlRepo struct {
	// dao is a collection of *sql.DB identified by strings (usernames), the point
	// being that each user gets their own connection pool
	dao  map[string]*sql.DB
	ddl  goyesql.Queries
	dml  goyesql.Queries
	conn string
}

// AppDB is an identifier for a specific *sql.DB in our dao map
const AppDB = "_app"

func NewPsqlRepo(conn string) *PsqlRepo {
	repo := PsqlRepo{}
	repo.conn = os.Getenv("DATABASE_URL")

	if repo.conn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	repo.dao = make(map[string]*sql.DB)
	repo.ensureDb(AppDB)

	// Create all the Tables, Views if they do not exist
	repo.ddl = goyesql.MustParseFile("data/sql/ddl.sql")

	logExec(repo.dao[AppDB], (string)(repo.ddl["create-user-table"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-command-table"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-tag-table"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-invocation-table"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-invocationtag-table"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-commandhistory-view"]))
	logExec(repo.dao[AppDB], (string)(repo.ddl["create-timestamp-index"]))

	// Load all data-manipulation queries
	repo.dml = goyesql.MustParseFile("data/sql/queries.sql")

	log.Println("storage init completed")
	return &repo
}

func (r *PsqlRepo) ql(identifier string) string {
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

// ensureDb verifies that we have created a connection pool for the given user
func (r *PsqlRepo) ensureDb(user string) {
	_, exists := r.dao[user]

	if !exists {
		var err error
		r.dao[user], err = sql.Open("postgres", r.conn)
		r.dao[user].SetMaxOpenConns(5)
		r.dao[user].SetMaxIdleConns(1)

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
func (r *PsqlRepo) InsertInvocations(user string, invocs Invocations) (err error) {
	r.ensureDb(user)
	tx, err := r.dao[user].Begin()
	if err != nil {
		return
	}
	for _, inv := range invocs {
		err = r.invocationTx(tx, user, inv)
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
func (r *PsqlRepo) invocationTx(tx *sql.Tx, user string, inv Invocation) (err error) {
	var cmdid int
	err = tx.QueryRow(r.ql("get-commandid-by-command"), inv.Command).Scan(&cmdid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(r.ql("insert-command"), inv.Command).Scan(&cmdid)
		if err2 != nil {
			return err2
		}
	}

	var invid int
	err = tx.QueryRow(r.ql("insert-invocation"), user, cmdid, inv.ExitCode,
		inv.Timestamp, inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&invid)

	for _, tag := range inv.Tags {
		err = r.AddTag(tx, user, invid, tag)
	}
	return
}

// AddTag persists a Tag to an Invocation, as part of a transaction.
func (r *PsqlRepo) AddTag(tx *sql.Tx, user string, invid int, tag string) (err error) {
	var tagid int
	err = tx.QueryRow(r.ql("get-tagid-by-name"), tag).Scan(&tagid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(r.ql("insert-tag"), tag).Scan(&tagid)
		if err2 != nil {
			return err2
		}
	}

	_, err = tx.Exec(r.ql("insert-invocationtag"), tagid, invid)
	return
}

// GetInvocations returns the [pageSize] most recent Invocations for the given
// user
func (r *PsqlRepo) GetInvocations(user string, pageSize int) (result Invocations, err error) {
	r.ensureDb(user)
	rows, err := r.dao[user].Query(r.ql("get-invocations-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}

	return queryInvocations(rows, pageSize)
}

// GetCommands returns the [pageSize] most recent Commands for the given user
func (r *PsqlRepo) GetCommands(user string, pageSize int) (result Commands, err error) {
	r.ensureDb(user)
	rows, err := r.dao[user].Query(r.ql("get-commands-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}
	return queryCommands(rows, pageSize)
}
