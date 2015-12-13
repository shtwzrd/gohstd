package gohstd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	gohst "github.com/warreq/gohstd/src/gohstd/common"
	"io/ioutil"
	"log"
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
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var commands interface{}
	if verbose {
		commands, err = repo.GetInvocations(user, pageSize)
	} else {
		commands, err = repo.GetCommands(user, pageSize)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(commands); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func CommandCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]

	var inv gohst.Invocations
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "plaintext; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprint(err))
		log.Println(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.Header().Set("Content-Type", "plaintext; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprint(err))
		return
	}
	if err := json.Unmarshal(body, &inv); err != nil {
		w.Header().Set("Content-Type", "plaintext; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Fprintf(w, fmt.Sprint(err))
		log.Println(err)
		return
	}

	err = repo.InsertInvocations(user, inv)
	if err != nil {
		w.Header().Set("Content-Type", "plaintext; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprint(err))
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func UserTagShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagCreate(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagRename(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
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
