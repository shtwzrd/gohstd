package service

import (
	"errors"
	"fmt"
	gohst "github.com/warreq/gohstd/common"
	bcrypt "golang.org/x/crypto/bcrypt"
	"net/http"
)

/*
* handlers are functions mapped to a route, which take *http.Request and a
* http.ResponseWriter. They are ultimately responsible for taking the Request
* and constructing the appropriate Response.
 */

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var u gohst.User
	err := ParseJsonEntity(r, &u)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	s, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	err = userRepo.InsertUser(u, s)
	if err != nil {
		if err.Error() == gohst.UserExistsError {
			HttpError(w, http.StatusConflict,
				errors.New(fmt.Sprintf("%s: %s", err, u.Username)))
			return
		}
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagRename(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}
