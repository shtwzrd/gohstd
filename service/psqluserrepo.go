package service

import (
	_ "github.com/lib/pq"
	gohst "github.com/warreq/gohstd/common"
)

// PsqlUserRepo is an implementation of a gohst UserRepo that uses
// PostgreSQL as a backing store
type PsqlUserRepo struct {
	dao PsqlDao
}

func NewPsqlUserRepo(psqldao PsqlDao) *PsqlUserRepo {
	return &PsqlUserRepo{psqldao}
}

func (r PsqlUserRepo) InsertUser(user gohst.User) (err error) {
	panic("Not yet implemented!")
}

func (r PsqlUserRepo) GetUserByName(uname string) (user gohst.User, err error) {
	panic("Not yet implemented!")
}
