package main

import (
	"log"
	"net/http"
)

// SessionUser returns the logged in user of a session or nil if no user is
// logged in.
func (s *server) SessionUser(r *http.Request) *User {
	session, err := s.cookies.Get(r, "sodago-session")
	if err != nil {
		log.Println("ERROR: failed to get session", session)
		return nil
	}
	name, ok := session.Values["user"].(string)
	if !ok || name == "" {
		return nil
	}
	return s.config.GetUser(name)
}
