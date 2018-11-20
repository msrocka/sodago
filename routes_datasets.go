package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterDataSetRoutes registers the data set methods to the given router.
func RegisterDataSetRoutes(r *mux.Router) {
	r.HandleFunc("/resource/{path}/{id}", GetDataSet)
	r.HandleFunc("/resource/{path}", PostDataSet).Methods("POST")
	r.HandleFunc("/resource/datastocks/{datastock}/{path}", GetDataSets).Methods("GET")
	r.HandleFunc("/resource/{path}", GetDataSets).Methods("GET")
}

// GetDataSet implements the `GET Dataset` function of the soda4LCA service API
func GetDataSet(w http.ResponseWriter, r *http.Request) {
	stock := db.RootDataStock()
	vars := mux.Vars(r)
	switch path := vars["path"]; path {
	case "processes":
		getProcess(vars, stock, w)
	case "flows":
		getFlow(vars, stock, w)
	case "flowproperties":
		getFlowProperty(vars, stock, w)
	case "unitgroups":
		getUnitGroup(vars, stock, w)
	case "contacts":
		getContact(vars, stock, w)
	case "sources":
		getSource(vars, stock, w)
	default:
		http.Error(w, "Unknown path "+path, http.StatusBadRequest)
	}
}

// GetDataSets implements the `GET Datasets` function of the soda4LCA service API
func GetDataSets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var stock *DataStock
	if stockID, found := vars["datastock"]; found {
		stock = db.DataStock(stockID)
	} else {
		stock = db.RootDataStock()
	}
	if stock == nil {
		http.Error(w, "Unknown data stock "+vars["datastock"],
			http.StatusBadRequest)
		return
	}

	content := db.Content(stock)
	list := &InfoList{}
	// TODO: filter by name; only return current version etc.
	switch path := vars["path"]; path {
	case "processes":
		list.Processes = content.Processes
	case "flows":
		list.Flows = content.Flows
	case "flowproperties":
		list.FlowProperties = content.FlowProperties
	case "unitgroups":
		list.UnitGroups = content.UnitGroups
	case "contacts":
		list.Contacts = content.Contacts
	case "sources":
		list.Sources = content.Sources
	default:
		http.Error(w, "Unknown path "+path, http.StatusBadRequest)
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

	var stock *DataStock
	stockID := r.Header.Get("stock")
	if stockID == "" {
		stock = db.RootDataStock()
	} else {
		stock = db.DataStock(stockID)
	}
	if stock == nil {
		http.Error(w, "Unknown data stock "+stockID, http.StatusBadRequest)
		return
	}

	switch path := mux.Vars(r)["path"]; path {
	case "processes":
		postProcess(data, stock, w)
	case "flows":
		postFlow(data, stock, w)
	case "flowproperties":
		postFlowProperty(data, stock, w)
	case "unitgroups":
		postUnitGroup(data, stock, w)
	case "contacts":
		postContact(data, stock, w)
	case "sources":
		postSource(data, stock, w)
	default:
		http.Error(w, "Unknown path "+path, http.StatusBadRequest)
	}
}
