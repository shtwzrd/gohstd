package main

import (
	"encoding/json"
	"fmt"
	"github.com/warreq/webgohst/Godeps/_workspace/src/github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

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
		commands, err = GetInvocations(user, pageSize)
	} else {
		commands, err = GetCommands(user, pageSize)
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

	var inv Invocations
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10485760))
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

	err = InsertInvocations(user, inv)
	if err != nil {
		w.Header().Set("Content-Type", "plaintext; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprint(err))
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func CommandTagCreate(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
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

func CommandTagDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func CommandDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func NotificationIndex(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}
