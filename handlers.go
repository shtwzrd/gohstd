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

	if err := json.NewEncoder(w).Encode(commands); err != nil {
		panic(err)
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
