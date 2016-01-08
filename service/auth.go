package service

import (
	"errors"
	"github.com/gorilla/mux"
	g "github.com/warreq/gohstd/common"
	bcrypt "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// AuthorizationFilter is a function type used for determining if a user is
// authorized to perform a given request. It returns true if a user meets the
// requirements of the authorization logic, and false otherwise.
type AuthorizationFilter func(user string, r *http.Request) bool

// AuthenticationBlackList contains the sets of service/Routes for which
// authentication is not necessary. Routes listed here bypass the authentication
// middleware and subsequently the authorization middleware as well, as to allow
// anonymous access to those endpoints.
var AuthenticationBlackList map[string]struct{}

// AuthorizationFilters contains the slice of functions for each service/Route
// to use for authorizing that endpoint. If a Route isn't found in the map, it
// requires no authorization
var AuthorizationFilters map[string][]AuthorizationFilter

func init() {
	AuthenticationBlackList = map[string]struct{}{
		"UserRegister": {},
		"GetPosts":     {},
	}

	AuthorizationFilters = map[string][]AuthorizationFilter{
		"CommandIndex":     []AuthorizationFilter{RequesterOwnsResource},
		"CommandCreate":    []AuthorizationFilter{RequesterOwnsResource},
		"CommandDelete":    []AuthorizationFilter{RequesterOwnsResource},
		"CommandTagCreate": []AuthorizationFilter{RequesterOwnsResource},
		"UserTagShow":      []AuthorizationFilter{RequesterOwnsResource},
		"UserTagRename":    []AuthorizationFilter{RequesterOwnsResource},
		"UserTagDelete":    []AuthorizationFilter{RequesterOwnsResource},
	}
}

// Auth is the http.Handler middleware for performinag authentication and
// authorization on a request before furthering the request to its actual handler.
// It consults the AuthenticationBlackList and AuthorizationBlackList to know
// whether to apply one, both or neither of these steps. Auth may short-circuit
// a request, returning http status codes 401 or 403. Otherwise, the request
//is handled as normal by the inner handler.
func Auth(inner http.Handler, route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, blacklisted := AuthenticationBlackList[route]
		if blacklisted {
			// Authentication isn't required -- jump directly to inner handler
			inner.ServeHTTP(w, r)
			return
		}
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Println(r)
			log.Println(username)
			log.Println(password)
			HttpError(w, http.StatusUnauthorized, errors.New(g.InvalidBasicAuthError))
			return
		}
		authenticated := Authenticate(username, password)
		if !authenticated {
			HttpError(w, http.StatusUnauthorized, nil)
			return
		}
		authorized := Authorize(username, r, route)
		if !authorized {
			HttpError(w, http.StatusForbidden, nil)
			return
		}
		// Success
		inner.ServeHTTP(w, r)
	})
}

// Authenticate compares the user and password provided and compares the stored
// credential for the given user against the given password. Authentication
// returns true in the case that the given password matches the stored
// credential and false otherwise.
func Authenticate(user, password string) bool {
	identity, err := userRepo.GetUserByName(user)
	if err != nil {
		log.Println(err)
		return false
	}
	credential := identity.Password
	err = bcrypt.CompareHashAndPassword([]byte(credential), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

// Authorize searches the AuthorizationFilters for functions pertaining to the
// given request and applies each filter in order. If any filter returns false,
// Authorize returns false. Otherwise and including the case of no filters,
// Authorize returns true.
func Authorize(user string, r *http.Request, route string) bool {
	for _, filter := range AuthorizationFilters[route] {
		if !filter(user, r) {
			return false
		}
	}
	return true
}

func RequesterOwnsResource(user string, r *http.Request) bool {
	vars := mux.Vars(r)
	owner := vars["username"]

	return user == owner
}
