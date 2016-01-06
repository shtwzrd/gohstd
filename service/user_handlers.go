package service

import (
	"errors"
	"fmt"
	gohst "github.com/warreq/gohstd/common"
	bcrypt "golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var u gohst.User
	err := ParseJsonEntity(r, &u)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	err2 := u.Validate()
	if err2 != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	s, err3 := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err3 != nil {
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	err4 := userRepo.InsertUser(u, s)
	if err4 != nil {
		if strings.Contains(err.Error(), gohst.UserExistsError) {
			HttpError(w, http.StatusConflict,
				errors.New(fmt.Sprintf("%s: %s", err, u.Username)))
			return
		}
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UserLogin is a do-nothing endpoint that merely confirms that authentication
// was successful. To that end, it is necessary that the endpoint is wrapped in
// actual authentication middleware.
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
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
