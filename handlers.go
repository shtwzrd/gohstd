package main

import (
	"encoding/json"
	"fmt"
	"github.com/warreq/webgohst/Godeps/_workspace/src/github.com/gorilla/mux"
	"log"
	"net/http"
)

func CommandIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	verbose := r.URL.Query().Get("verbose") == "true"

	if verbose {
		commands, err := GetAllInvocations(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if err := json.NewEncoder(w).Encode(commands); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		commands, err := GetAllCommands(user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if err := json.NewEncoder(w).Encode(commands); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CommandShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commandId := vars["commandId"]
	fmt.Fprintln(w, "Command show: ", commandId)
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func CommandCreate(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
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
