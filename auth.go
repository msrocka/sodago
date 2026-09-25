package main

import (
	"net/http"
	"strings"
)

// bearerPrefix is the prefix of a token in the Authorization header.
const bearerPrefix = "Bearer "

// permissionDenied writes the response of the server when a request contains
// invalid credentials or is not allowed (similar to soda4LCA).
func permissionDenied(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Permission denied."))
}

// bearerToken returns the token of the Authorization header or an empty string
// when no token is sent.
func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(header) < len(bearerPrefix) ||
		!strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(header[len(bearerPrefix):])
}

// TokenUser returns the user that is authenticated with the given token.
func (s *server) TokenUser(token string) (*User, error) {
	// a token that is configured for a user is accepted directly
	if user := s.config.GetUserByToken(token); user != nil {
		return user, nil
	}
	name, err := s.tokens.Verify(token)
	if err != nil {
		return nil, err
	}
	user := s.config.GetUser(name)
	if user == nil {
		return nil, errInvalidToken
	}
	return user, nil
}

// RequestUser returns the user that is authenticated with the given request.
// This can be a user from a session cookie or from an authentication token. It
// returns nil when the request does not contain any credentials and an error
// when it contains invalid credentials.
func (s *server) RequestUser(r *http.Request) (*User, error) {
	if token := bearerToken(r); token != "" {
		return s.TokenUser(token)
	}
	return s.SessionUser(r), nil
}

// authenticationMiddleware checks the credentials of a request before it is
// handled: reading is also allowed for anonymous users, while writing requires
// a valid session or token. Requests with invalid credentials are rejected.
func (s *server) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// the authentication routes are not protected
		if strings.HasPrefix(r.URL.Path, authenticatePath) {
			next.ServeHTTP(w, r)
			return
		}
		user, err := s.RequestUser(r)
		if err != nil || (user == nil && isWrite(r)) {
			permissionDenied(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// isWrite returns true when the request modifies data on the server.
func isWrite(r *http.Request) bool {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}
