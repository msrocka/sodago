package main

import (
	"net/http"
)

func getProcess(params map[string]string, stock *DataStock, w http.ResponseWriter) {
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
	var selected *ProcessInfo
	var selectedVersion *Version
	all := db.Content(stock).Processes
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

func getFlowProperty(params map[string]string, stock *DataStock, w http.ResponseWriter) {
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
	var selected *FlowPropertyInfo
	var selectedVersion *Version
	all := db.Content(stock).FlowProperties
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

func getUnitGroup(params map[string]string, stock *DataStock, w http.ResponseWriter) {
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
	var selected *UnitGroupInfo
	var selectedVersion *Version
	all := db.Content(stock).UnitGroups
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

func getContact(params map[string]string, stock *DataStock, w http.ResponseWriter) {
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
	var selected *ContactInfo
	var selectedVersion *Version
	all := db.Content(stock).Contacts
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

func getSource(params map[string]string, stock *DataStock, w http.ResponseWriter) {
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
	var selected *SourceInfo
	var selectedVersion *Version
	all := db.Content(stock).Sources
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
