package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (s *server) handleGetLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := strings.TrimSpace(vars["user"])
		pw := strings.TrimSpace(vars["password"])

		// find the user
		if name == "" || pw == "" {
			http.Error(w, "user name and password must have a value",
				http.StatusBadRequest)
			return
		}
		currentUser := s.SessionUser(r)
		if currentUser != nil && currentUser.Name == name {
			w.Write([]byte("You are already logged in as a user"))
			return
		}
		user := s.config.GetUser(name)
		if user == nil {
			http.Error(w, "incorrect password or user name",
				http.StatusBadRequest)
			return
		}

		// check the password (stored as plain text, testing only)
		if user.Password != pw {
			http.Error(w, "incorrect password or user name",
				http.StatusBadRequest)
			return
		}

		session, err := s.cookies.Get(r, "sodago-session")
		if err != nil {
			log.Println("ERROR: could not get session", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		session.Values["user"] = user.Name
		session.Save(r, w)
		w.Write([]byte("Login successful"))
	}
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
		user := s.SessionUser(r)
		if user != nil {
			response.IsAuthenticated = true
			response.UserName = user.Name
			response.Roles = user.Roles
		}
		writeXML(response, w)
	}
}
