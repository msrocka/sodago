package main

import (
	"net/http"
)

func postProcess(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadProcessInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsProcess(info) {
		http.Error(w, "Process "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.Processes = append(content.Processes, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}

func postFlow(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadFlowInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsFlow(info) {
		http.Error(w, "Flow "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.Flows = append(content.Flows, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}

func postFlowProperty(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadFlowPropertyInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsFlowProperty(info) {
		http.Error(w, "FlowProperty "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.FlowProperties = append(content.FlowProperties, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}

func postUnitGroup(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadUnitGroupInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsUnitGroup(info) {
		http.Error(w, "UnitGroup "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.UnitGroups = append(content.UnitGroups, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}

func postContact(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadContactInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsContact(info) {
		http.Error(w, "Contact "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.Contacts = append(content.Contacts, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}

func postSource(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadSourceInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsSource(info) {
		http.Error(w, "Source "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	db.Put(DataSetBucket, info.Key(stock), data)
	content.Sources = append(content.Sources, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}
