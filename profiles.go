package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Profile contains information of an EPD profile.
type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetProfileDescriptors returns a JSON array with the EPD profile descriptors
// from the server.
func GetProfileDescriptors(w http.ResponseWriter, r *http.Request) {
	dir := "data/profiles"
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		ServeJSONBytes([]byte("[]"), w)
		return
	}
	files, err := ioutil.ReadDir("data/profiles")
	if err != nil {
		ServeJSONBytes([]byte("[]"), w)
		return
	}
	profiles := make([]*Profile, 0)
	for _, file := range files {
		data, err := ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			continue
		}
		p := &Profile{}
		err = json.Unmarshal(data, p)
		if err != nil {
			continue
		}
		profiles = append(profiles, p)
	}
	ServeJSON(profiles, w)
}

// GetProfile returns an EPD profile for a given ID.
func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	path := "data/profiles/" + id + ".json"
	if _, err := os.Stat(path); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ServeJSONBytes(data, w)
}
