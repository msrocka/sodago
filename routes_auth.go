package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes registers the authentication and authorization routes
// methods to the given router.
func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/resource/authenticate/login", Login)
}

// Login implements the `Login` function of the soda4LCA service API
func Login(w http.ResponseWriter, r *http.Request) {
	// TODO checks nothing currently
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Logged in"))
}
