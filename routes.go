package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) registerRoutes(r *mux.Router, args *Args) {

	// data stocks
	r.HandleFunc("/resource/datastocks", s.handleGetDataStocks()).Methods("GET")

	// profiles
	r.Methods("GET").Path("/resource/profiles").
		HandlerFunc(GetProfileDescriptors)
	r.Methods("GET").Path("/resource/profiles/").
		HandlerFunc(GetProfileDescriptors)
	r.Methods("GET").Path("/resource/profiles/{id}").
		HandlerFunc(GetProfile)

	// login: currently nothing is checked here
	r.HandleFunc("/resource/authenticate/login",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Logged in"))
		})

	// GET a single data set from the root data stock
	// specific version
	r.HandleFunc("/resource/{path}/{id}", s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/{path}/{id}", s.handleGetDataSet()).
		Methods("GET", "HEAD")

	// GET a single data set from a data stock
	// specific version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).Methods("GET", "HEAD")

	// POST a data set
	r.HandleFunc("/resource/{path}", s.handlePostDataSet()).Methods("POST")

	// GET  a list of data sets
	r.HandleFunc("/resource/datastocks/{datastock}/{path}", GetDataSets).Methods("GET")
	r.HandleFunc("/resource/{path}", GetDataSets).Methods("GET")
}
