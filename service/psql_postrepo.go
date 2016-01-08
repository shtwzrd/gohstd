package service

import (
	"errors"
	_ "github.com/lib/pq"
	g "github.com/warreq/gohstd/common"
	"time"
)

// PsqlPostRepo is an implementation of a gohst PostRepo that uses
// PostgreSQL as a backing store
type PsqlPostRepo struct {
	dao PsqlDao
}

func NewPsqlPostRepo(dao PsqlDao) *PsqlPostRepo {
	return &PsqlPostRepo{dao}
}

func (r PsqlPostRepo) InsertPost(post g.NewPost, username string) error {
	db := r.dao.EnsurePool(username)
	_, err := db.Exec(r.dao.Query("insert-post"),
		username, post.Title, post.Message, time.Now())
	if err != nil {
		return errors.New("Could not persist entity")
	}
	return nil
}

func (r PsqlPostRepo) DeletePost(user string, postid int) error {
	db := r.dao.EnsurePool(user)
	_, err := db.Exec(r.dao.Query("delete-post"), postid)
	return err
}

func (r PsqlPostRepo) GetPosts() (posts g.Posts, err error) {
	db := r.dao.EnsurePool(AppDB)
	rows, err := db.Query(r.dao.Query("get-posts"))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tmp g.Post
		err = rows.Scan(&tmp.Username, &tmp.Title, &tmp.Message, &tmp.Timestamp)

		if err != nil {
			return nil, err
		}

		posts = append(posts, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}
