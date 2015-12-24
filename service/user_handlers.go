package service

import (
	// "encoding/json"
	// "fmt"
	// "github.com/gorilla/mux"
	gohst "github.com/warreq/gohstd/common"
	// "io/ioutil"
	// "log"
	"net/http"
	// "strconv"
)

/*
* handlers are functions mapped to a route, which take *http.Request s and a
* http.ResponseWriter. They are ultimately responsible for taking the Request
* and constructing the appropriate Response.
 */

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var user gohst.User
	err := ParseJsonEntity(r, &user)

	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}
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
