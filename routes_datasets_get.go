package main

import (
	"net/http"
)

func getFlow(params map[string]string, stock *DataStock, w http.ResponseWriter) {
	id, ok := params["id"]
	if !ok {
		http.Error(w, "No ID provided", http.StatusBadRequest)
		return
	}
	v, ok := params["version"]
	var version *Version
	if ok {
		version = ParseVersion(v)
	}
	var selected *FlowInfo
	var selectedVersion *Version
	all := db.Content(stock).Flows
	for i := range all {
		info := &all[i]
		if !Matches(info, id, version) {
			continue
		}
		if version != nil {
			selected = info
			break
		}
		infoV := ParseVersion(info.Version)
		if selected == nil || infoV.NewerThan(selectedVersion) {
			selected = info
			selectedVersion = infoV
		}
	}
	ServeDataSet(selected, stock, w)
}
