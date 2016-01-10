package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	gohst "github.com/warreq/gohstd/common"
	"net/http"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postRepo.GetPosts()

	if err != nil {
		HttpError(w, http.StatusBadRequest, nil)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		HttpError(w, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SendPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	var p gohst.NewPost
	err := ParseJsonEntity(r, &p)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	err = postRepo.InsertPost(p, user)
	if err != nil {
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
