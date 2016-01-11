package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	g "github.com/warreq/gohstd/common"
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

func InvocationToBase64(invoc g.Invocation) g.Invocation {
	enc := base64.StdEncoding
	invoc.Command = enc.EncodeToString([]byte(invoc.Command))
	invoc.User = enc.EncodeToString([]byte(invoc.User))
	invoc.Host = enc.EncodeToString([]byte(invoc.Host))
	invoc.Shell = enc.EncodeToString([]byte(invoc.Shell))
	invoc.Directory = enc.EncodeToString([]byte(invoc.Directory))
	for i, t := range invoc.Tags {
		invoc.Tags[i] = enc.EncodeToString([]byte(t))
	}
	return invoc
}

func InvocationsToBase64(invocs g.Invocations) g.Invocations {
	for i, v := range invocs {
		invocs[i] = InvocationToBase64(v)
	}
	return invocs
}

func InvocationFromBase64(invoc g.Invocation) (g.Invocation, error) {
	enc := base64.StdEncoding
	var err error
	cmd, err := enc.DecodeString(invoc.Command)
	user, err := enc.DecodeString(invoc.User)
	host, err := enc.DecodeString(invoc.Host)
	shell, err := enc.DecodeString(invoc.Shell)
	dir, err := enc.DecodeString(invoc.Directory)
	tags := make([][]byte, len(invoc.Tags))
	for i, t := range invoc.Tags {
		tags[i], err = enc.DecodeString(t)
	}
	if err != nil {
		return invoc, err
	}
	invoc.Command = string(cmd)
	invoc.User = string(user)
	invoc.Host = string(host)
	invoc.Shell = string(shell)
	invoc.Directory = string(dir)
	for i, t := range tags {
		invoc.Tags[i] = string(t)
	}
	return invoc, nil
}

func InvocationsFromBase64(invocs g.Invocations) (g.Invocations, error) {
	var err error
	for i, v := range invocs {
		invocs[i], err = InvocationFromBase64(v)
	}
	return invocs, err
}

func CommandFromBase64(cmd g.Command) (g.Command, error) {
	enc := base64.StdEncoding
	c, err := enc.DecodeString(string(cmd))
	if err != nil {
		return cmd, nil
	}
	return g.Command(c), nil
}

func CommandsFromBase64(cmds g.Commands) (g.Commands, error) {
	var err error
	for i, v := range cmds {
		cmds[i], err = CommandFromBase64(v)
	}
	return cmds, err
}
