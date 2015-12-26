package service

import (
	"errors"
	_ "github.com/lib/pq"
	g "github.com/warreq/gohstd/common"
	"strings"
)

// PsqlUserRepo is an implementation of a gohst UserRepo that uses
// PostgreSQL as a backing store
type PsqlUserRepo struct {
	dao PsqlDao
}

func NewPsqlUserRepo(dao PsqlDao) *PsqlUserRepo {
	return &PsqlUserRepo{dao}
}

func (r PsqlUserRepo) InsertUser(user g.User, secret g.Secret) (err error) {
	db := r.dao.EnsurePool(user.Username)
	_, err = db.Exec(r.dao.Query("insert-user"),
		user.Username, user.Email, string(secret))
	if strings.Contains(err.Error(), `duplicate key value`) {
		return errors.New(g.UserExistsError)
	}

	return err
}

func (r PsqlUserRepo) GetUserByName(uname string) (user g.User, err error) {
	db := r.dao.EnsurePool(uname)
	var email string
	err = db.QueryRow(r.dao.Query("get-user-by-name"), uname).Scan(&email)

	user = g.User{}
	user.Username = uname
	user.Email = email

	return
}
