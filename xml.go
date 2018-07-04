package main

import (
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
