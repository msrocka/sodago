package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterDataSetRoutes registers the data set methods to the given router.
func RegisterDataSetRoutes(r *mux.Router) {
	r.HandleFunc("/resource/{path}", PostDataSet).Methods("POST")
}

// PostDataSet handles a post request of a data set.
func PostDataSet(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	info := ReadFlowInfo(data)
	if info != nil {
		ServeXML(info, w)
	}

	println(path)
	println(string(data))
}
