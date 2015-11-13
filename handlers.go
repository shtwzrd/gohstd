package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func CommandIndex(w http.ResponseWriter, r *http.Request) {
	commands := Commands{
		Command{Id: 1},
		Command{Id: 2},
	}

	if err := json.NewEncoder(w).Encode(commands); err!= nil {
		panic(err)
	}
}

func CommandShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commandId := vars["commandId"]
	fmt.Fprintln(w, "Command show: ", commandId)
}