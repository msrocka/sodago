package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func writeXML(e interface{}, w http.ResponseWriter) {
	if e == nil {
		http.Error(w, "No data", http.StatusInternalServerError)
		return
	}
	data, err := xml.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBytesXML(data, w)
}

func writeBytesXML(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write(data)
}

func writeJSON(e interface{}, w http.ResponseWriter) {
	if e == nil {
		http.Error(w, "No data", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBytesJSON(data, w)
}

func writeBytesJSON(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
