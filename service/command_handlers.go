package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	gohst "github.com/warreq/gohstd/common"
	"net/http"
	"strconv"
)

/*
* handlers are functions mapped to a route, which take *http.Request s and a
* http.ResponseWriter. They are ultimately responsible for taking the Request
* and constructing the appropriate Response.
 */

// CommandIndex is the handler for querying a user's commands, along with any
// query parameters
func CommandIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	count := r.URL.Query().Get("count")
	verbose := r.URL.Query().Get("verbose") == "true"

	var pageSize int
	var err error
	if count == "" {
		pageSize = 0
	} else {
		pageSize, err = strconv.Atoi(count)
		if err != nil {
			HttpError(w, http.StatusBadRequest, nil)
			return
		}
	}

	var commands interface{}
	if verbose {
		commands, err = commandRepo.GetInvocations(user, pageSize)
	} else {
		commands, err = commandRepo.GetCommands(user, pageSize)
	}
	if err != nil {
		HttpError(w, http.StatusBadRequest, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(commands); err != nil {
		HttpError(w, http.StatusInternalServerError, nil)
	}
}

func CommandCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]

	var inv gohst.Invocations
	err := ParseJsonEntity(r, &inv)

	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	err = commandRepo.InsertInvocations(user, inv)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func CommandDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func CommandTagCreate(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func CommandTagDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}
