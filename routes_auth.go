package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// handlePostLogin handles the session login with the user name and password
// given as form parameters.
func (s *server) handlePostLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "could not read form", http.StatusBadRequest)
			return
		}
		s.login(w, r, r.FormValue("username"), r.FormValue("password"))
	}
}

// handleGetLogin handles the legacy session login with the user name and
// password given as query parameters.
func (s *server) handleGetLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		s.login(w, r, vars["user"], vars["password"])
	}
}

// login creates a session for the user with the given credentials.
func (s *server) login(w http.ResponseWriter, r *http.Request, name, pw string) {
	name = strings.TrimSpace(name)
	pw = strings.TrimSpace(pw)
	if name == "" || pw == "" {
		http.Error(w, "user name and password must have a value",
			http.StatusUnauthorized)
		return
	}

	currentUser := s.SessionUser(r)
	if currentUser != nil && currentUser.Name == name {
		w.Write([]byte("You are already logged in as a user"))
		return
	}

	user := s.checkPassword(name, pw)
	if user == nil {
		http.Error(w, "incorrect password or user name",
			http.StatusUnauthorized)
		return
	}

	session, err := s.cookies.Get(r, "sodago-session")
	if err != nil {
		log.Println("ERROR: could not get session", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	session.Values["user"] = user.Name
	if err := session.Save(r, w); err != nil {
		log.Println("ERROR: could not save session", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Login successful"))
}

// handlePostToken handles the token request with the user name and password
// given as form parameters.
func (s *server) handlePostToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "could not read form", http.StatusBadRequest)
			return
		}
		s.writeToken(w, r.FormValue("username"), r.FormValue("password"))
	}
}

// handleGetToken handles the legacy token request with the user name and
// password given as query parameters.
func (s *server) handleGetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		s.writeToken(w, vars["user"], vars["password"])
	}
}

// writeToken creates a token for the user with the given credentials.
func (s *server) writeToken(w http.ResponseWriter, name, pw string) {
	name = strings.TrimSpace(name)
	pw = strings.TrimSpace(pw)
	if name == "" || pw == "" {
		http.Error(w, "user name and password must have a value",
			http.StatusUnauthorized)
		return
	}
	user := s.checkPassword(name, pw)
	if user == nil {
		http.Error(w, "permission denied", http.StatusUnauthorized)
		return
	}
	token, err := s.tokens.Generate(user, time.Now())
	if err != nil {
		log.Println("ERROR: could not create token", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(token))
}

// checkPassword returns the user with the given name when the password is
// correct. Note: the passwords are stored as plain text (testing only).
func (s *server) checkPassword(name, pw string) *User {
	user := s.config.GetUser(name)
	if user == nil || user.Password != pw {
		return nil
	}
	return user
}

// handleGetLogout closes the session of the currently logged in user.
func (s *server) handleGetLogout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if s.SessionUser(r) == nil {
			w.Write([]byte("currently not authenticated"))
			return
		}

		session, err := s.cookies.Get(r, "sodago-session")
		if err != nil {
			log.Println("ERROR: could not get session", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		// remove the user from the session and delete the session cookie
		delete(session.Values, "user")
		session.Options.MaxAge = -1
		if err := session.Save(r, w); err != nil {
			log.Println("ERROR: could not save session", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("successfully logged out"))
	}
}

func (s *server) handleGetAuthenticationStatus() http.HandlerFunc {

	type response struct {
		XMLName         xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI authInfo"`
		IsAuthenticated bool     `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI authenticated"`
		UserName        string   `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI userName,omitempty"`
		Roles           []string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := &response{}
		// invalid credentials are ignored here as this route is not protected
		if user, err := s.RequestUser(r); err == nil && user != nil {
			response.IsAuthenticated = true
			response.UserName = user.Name
			response.Roles = user.Roles
		}
		writeXML(response, w)
	}
}
