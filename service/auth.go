package service

import (
	"errors"
	g "github.com/warreq/gohstd/common"
	bcrypt "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// AuthenticationBlackList contains the sets of serivce/Routes for which
// authentication is not necessary. Routes listed here bypass the authentication
// middleware and subsequently the authorization middleware as well, as to allow
// anonymous access to those endpoints.
var AuthenticationBlackList map[string]struct{}

// AuthorizationBlackList contains the sets of serivce/Routes for which
// authorization is not necessary, but authentication is required.
var AuthorizationBlackList map[string]struct{}

func init() {
	AuthenticationBlackList = map[string]struct{}{
		"UserRegister": {},
	}

	AuthorizationBlackList = map[string]struct{}{}
}

// Auth is the http.Handler middleware for performinag authentication and
// authorization on a request before furthering the request to its actual handler.
// It consults the AuthenticationBlackList and AuthorizationBlackList to know
// whether to apply one, both or neither of these steps. Auth may short-circuit
// a request, returning http status codes 400, 401 or 403. Otherwise, the request
//is handled as normal by the inner handler.
func Auth(inner http.Handler, route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, blacklisted := AuthenticationBlackList[route]
		if blacklisted {
			// Authentication isn't required -- jump directly to inner handler
			inner.ServeHTTP(w, r)
		}
		username, password, ok := r.BasicAuth()
		if !ok {
			HttpError(w, http.StatusBadRequest, errors.New(g.InvalidBasicAuthError))
			return
		}
		authenticated := Authenticate(username, password, route)
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

func Authenticate(user, password, route string) bool {
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

func Authorize(user string, r *http.Request, route string) bool {
	return true
}
