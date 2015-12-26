package service

import (
	"errors"
	g "github.com/warreq/gohstd/common"
	"net/http"
)

// AuthenticationBlackList contains the sets of serivce/Routes for which
// authentication is not necessary. Routes listed here bypass the authentication
// middleware and subsequently the authorization middleware as well, as to allow
// anonymous access to those endpoints.
var AuthenicationBlackList map[string]struct{}

func init() {
	AuthenicationBlackList = map[string]struct{}{
		"UserRegister": {},
	}
}

// Auth is the http.Handler middleware for performinag authentication and
// authorization on a request before furthering the request to its actual handler.
// It consults the AuthenticationBlackList and AuthorizationBlackList to know
// whether to apply one, both or neither of these steps. Auth may short-circuit
// a request, returning http status codes 400, 401 or 403. Otherwise, the request
//is handled as normal by the inner handler.
func Auth(inner http.Handler, route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			HttpError(w, http.StatusBadRequest, errors.New(g.InvalidBasicAuthError))
			return
		}
		authenticated := Authenticate(username, password)
		if !authenticated {
			HttpError(w, http.StatusUnauthorized, nil)
			return
		}
		authorized := Authorize(username, route)
		if !authorized {
			HttpError(w, http.StatusForbidden, nil)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

func Authenticate(user, password string) bool {
	return true
}

func Authorize(user, route string) bool {
	return true
}
