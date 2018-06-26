package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// profiles
	router.Methods("GET").Path("/resource/profiles").
		HandlerFunc(GetProfileDescriptors)
	router.Methods("GET").Path("/resource/profiles/").
		HandlerFunc(GetProfileDescriptors)
	router.Methods("GET").Path("/resource/profiles/{id}").
		HandlerFunc(GetProfile)
	http.ListenAndServe(":8080", router)
}
