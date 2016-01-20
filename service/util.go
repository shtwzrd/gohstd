package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HttpError is a convenience function for writing the necessary headers and
// content for returning an error
func HttpError(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	w.Header().Set("Content-Type", "plaintext;charset=UTF-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, fmt.Sprint(err))
}

func ParseJsonEntity(r *http.Request, entity interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, &entity); err != nil {
		return err
	}
	return nil
}
