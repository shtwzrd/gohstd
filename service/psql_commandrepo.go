package service

import (
	"database/sql"
	_ "github.com/lib/pq"
	gohst "github.com/warreq/gohstd/common"
	"log"
	"strings"
)

// PsqlCommandRepo is an implementation of a gohst CommandRepo that uses
// PostgreSQL as a backing store
type PsqlCommandRepo struct {
	dao PsqlDao
}

func NewPsqlCommandRepo(psqldao PsqlDao) *PsqlCommandRepo {
	return &PsqlCommandRepo{psqldao}
}

// queryCommands is a common handler for implementing paging over Commands
func queryCommands(rows *sql.Rows, pageSize int) (result gohst.Commands, err error) {
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
		result = append(result, gohst.Command(c))
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
func queryInvocations(rows *sql.Rows, pageSize int) (result gohst.Invocations, err error) {
	defer rows.Close()
	var tmp gohst.Invocation
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
func (r PsqlCommandRepo) InsertInvocations(user string, invocs gohst.Invocations) (err error) {
	db := r.dao.EnsurePool(user)
	tx, err := db.Begin()
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
func (r *PsqlCommandRepo) invocationTx(tx *sql.Tx, user string, inv gohst.Invocation) (err error) {
	var cmdid int
	err = tx.QueryRow(r.dao.Query("get-commandid-by-command"), inv.Command).
		Scan(&cmdid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(r.dao.Query("insert-command"), inv.Command).Scan(&cmdid)
		if err2 != nil {
			return err2
		}
	}

	var invid int
	err = tx.QueryRow(r.dao.Query("insert-invocation"), user, cmdid, inv.ExitCode,
		inv.Timestamp, inv.Host, inv.User, inv.Shell, inv.Directory).Scan(&invid)

	for _, tag := range inv.Tags {
		err = r.AddTag(tx, user, invid, tag)
	}
	return
}

// AddTag persists a Tag to an Invocation, as part of a transaction.
func (r PsqlCommandRepo) AddTag(tx *sql.Tx, user string, invid int, tag string) (err error) {
	var tagid int
	err = tx.QueryRow(r.dao.Query("get-tagid-by-name"), tag).Scan(&tagid)

	if err == sql.ErrNoRows {
		err2 := tx.QueryRow(r.dao.Query("insert-tag"), tag).Scan(&tagid)
		if err2 != nil {
			return err2
		}
	}

	_, err = tx.Exec(r.dao.Query("insert-invocationtag"), tagid, invid)
	return
}

// GetInvocations returns the [pageSize] most recent Invocations for the given
// user
func (r PsqlCommandRepo) GetInvocations(user string, pageSize int) (result gohst.Invocations, err error) {
	db := r.dao.EnsurePool(user)
	rows, err := db.Query(r.dao.Query("get-invocations-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}

	return queryInvocations(rows, pageSize)
}

// GetCommands returns the [pageSize] most recent Commands for the given user
func (r PsqlCommandRepo) GetCommands(user string, pageSize int) (result gohst.Commands, err error) {
	db := r.dao.EnsurePool(user)
	rows, err := db.Query(r.dao.Query("get-commands-by-user"), user, pageSize)
	if err != nil {
		log.Println(err)
		return
	}
	return queryCommands(rows, pageSize)
}
