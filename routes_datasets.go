package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
