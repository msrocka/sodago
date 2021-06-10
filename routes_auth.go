package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleGetLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["user"]
		pw := vars["password"]

		// find the user
		if name == "" || pw == "" {
			http.Error(w, "user name and password must have a value",
				http.StatusBadRequest)
			return
		}
		currentUser := s.SessionUser(r)
		if currentUser.Name == name {
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
