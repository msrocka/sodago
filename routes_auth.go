package main

import (
	"encoding/base64"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

		// check the password
		hash, err := base64.StdEncoding.DecodeString(user.Hash)
		if err != nil {
			log.Println("ERROR: could not decode user hash", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword(hash, []byte(pw))
		if err != nil {
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
