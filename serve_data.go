package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// ServeXML converts the given entity to a XML string and writes it to the
// given response.
func ServeXML(e interface{}, w http.ResponseWriter) {
	if e == nil {
		http.Error(w, "No data", http.StatusInternalServerError)
		return
	}
	data, err := xml.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ServeXMLBytes(data, w)
}

// ServeXMLBytes writes the given data as XML content to the given writer. It
// also sets the respective access control headers so that cross domain requests
// are supported.
func ServeXMLBytes(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write(data)
}

// ServeJSON converts the given entity to a JSON string and writes it to the
// given response.
func ServeJSON(e interface{}, w http.ResponseWriter) {
	if e == nil {
		http.Error(w, "No data", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ServeJSONBytes(data, w)
}

// ServeJSONBytes writes the given data as JSON content to the given writer. It
// also sets the respective access control headers so that cross domain requests
// are supported.
func ServeJSONBytes(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
