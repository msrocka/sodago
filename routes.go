package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router, args *Args) {
	log.Println("Register routes with static files from:", args.StaticDir)

	// data stocks
	r.HandleFunc("/resource/datastocks", GetDataStocks).Methods("GET")

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

	// GET a single data set
	r.HandleFunc("/resource/{path}/{id}", GetDataSet).
		Methods("GET", "HEAD")
	r.HandleFunc("/resource/{path}/{id}", GetDataSet).
		Queries("version", "{version}").Methods("GET", "HEAD")

	// POST a data set
	r.HandleFunc("/resource/{path}", PostDataSet).Methods("POST")

	// GET  a list of data sets
	r.HandleFunc("/resource/datastocks/{datastock}/{path}", GetDataSets).Methods("GET")
	r.HandleFunc("/resource/{path}", GetDataSets).Methods("GET")
}
