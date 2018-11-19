package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterDataSetRoutes registers the data set methods to the given router.
func RegisterDataSetRoutes(r *mux.Router) {
	r.HandleFunc("/resource/{path}", PostDataSet).Methods("POST")
	r.HandleFunc("/resource/{path}", GetDataSets).Methods("GET")
}

// GetDataSets implements the `GET Datasets` request of the soda4LCA service API
func GetDataSets(w http.ResponseWriter, r *http.Request) {
	stock := db.RootDataStock()
	content := db.Content(stock)
	list := &InfoList{}
	switch path := mux.Vars(r)["path"]; path {
	case "flows":
		list.Flows = content.Flows
		// TODO
	}
	ServeXML(list, w)
}

// PostDataSet handles a post request of a data set.
func PostDataSet(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read body "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	stock := db.RootDataStock()
	switch path := mux.Vars(r)["path"]; path {
	case "flows":
		postFlow(data, stock, w)
	default:
		http.Error(w, "Unknown path "+path, http.StatusBadRequest)
	}
}
