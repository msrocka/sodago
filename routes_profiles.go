package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (s *server) handleGetProfiles() http.HandlerFunc {

	type profile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dir := filepath.Join(s.dir.root, "profiles")
		if !fileExists(dir) {
			writeBytesJSON([]byte("[]"), w)
			return
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			http.Error(w, "server error: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		profiles := make([]*profile, 0)
		for _, file := range files {
			data, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				http.Error(w, "server error: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
			p := &profile{}
			if err := json.Unmarshal(data, p); err != nil {
				http.Error(w, "server error: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
			profiles = append(profiles, p)
		}
		writeJSON(profiles, w)
	}
}

func (s *server) handleGetProfile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		file := filepath.Join(s.dir.root, "profiles", id+".json")
		if !fileExists(file) {
			http.Error(w, "Profile "+id+"does not exist", http.StatusNotFound)
			return
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			http.Error(w, "server error: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		writeBytesJSON(data, w)
	}
}
