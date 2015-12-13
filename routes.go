package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Route represents a URI route for the server to support. A definition of a
// Route declares support for an HTTP endpoint and maps that URI to the
// appropriate HTTP handler.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var repo Repo

// NewRouter constructs a *mux.Router based on the routes defined in this
// package, which can then be passed to the net/http server.
func NewRouter(r Repo) *mux.Router {

	repo = r
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = StandardHeader(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"UserRegister",
		"POST",
		"/api/users/register",
		UserRegister,
	},
	Route{
		"UserShow",
		"GET",
		"/api/users/{username}",
		UserShow,
	},
	Route{
		"CommandIndex",
		"GET",
		"/api/users/{username}/commands",
		CommandIndex,
	},
	Route{
		"CommandCreate",
		"POST",
		"/api/users/{username}/commands",
		CommandCreate,
	},
	Route{
		"CommandTagCreate",
		"POST",
		"/api/users/{username}/commands/{commandId}/tags",
		CommandTagCreate,
	},
	Route{
		"UserTagShow",
		"GET",
		"/api/users/{username}/tags",
		UserTagShow,
	},
	Route{
		"UserTagRename",
		"PUT",
		"/api/users/{username}/tags/{tag}",
		UserTagRename,
	},
	Route{
		"UserTagDelete",
		"DELETE",
		"/api/users/{username}/tags/{tag}",
		UserTagDelete,
	},
	Route{
		"CommandTagDelete",
		"DELETE",
		"/api/users/{username}/commands/{commandId/tags/{tag}",
		CommandTagDelete,
	},
	Route{
		"CommandDelete",
		"DELETE",
		"/api/users/{username}/commands/{commandId}",
		CommandDelete,
	},
}
