package main

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *server) handleGetExternalFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		stockID := vars["datastock"]
		uid := vars["id"]
		file := vars["file"]

		stock := s.dir.findDataStock(stockID)
		if stock == nil {
			http.Error(w, "Unknown data stock", http.StatusNotFound)
			return
		}

		path := filepath.Join(stock.dir, "external_docs", uid, file)
		if !fileExists(path) {
			http.Error(w, "Unknown file", http.StatusNotFound)
			return
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			http.Error(w, "Could not read file", http.StatusInternalServerError)
			return
		}

		parts := strings.Split(file, ".")
		if len(parts) > 1 {
			ext := "." + strings.ToLower(parts[len(parts)-1])
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	}
}

func (s *server) handlePostSourceWithFiles() http.HandlerFunc {

	type response struct {
		IsRoot    bool   `xml:"root,attr"`
		ID        string `xml:"uuid"`
		ShortName string `xml:"shortName"`
		Name      string `xml:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "failed to parse multi-part form: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		first := func(vals []string) string {
			if len(vals) == 0 {
				return ""
			}
			return vals[0]
		}

		// lock the mutex until this method is ready
		s.mutex.Lock()
		defer s.mutex.Unlock()

		// try to read the source
		sourceXML := first(r.Form["file"])
		if sourceXML == "" {
			http.Error(w, "no source XML stored in `file`", http.StatusBadRequest)
			return
		}
		source := []byte(sourceXML)

		// extract the source info
		sourceInfo, err := extractIndexEntry(sourcePath, source)
		if err != nil {
			http.Error(w, "failed to read `file` param: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// try to save the source
		stockID := first(r.Form["stock"])
		stock, err := s.dir.put(stockID, sourcePath, source)
		if err != nil {
			http.Error(w, "failed to store source: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// create the folder for external documents of that source
		dir := filepath.Join(stock.dir, "external_docs", sourceInfo.UUID)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// save the uploaded files
		for f := range r.Form {
			if f == "file" {
				continue
			}
			str := first(r.Form[f])
			if str == "" {
				continue
			}
			path := filepath.Join(dir, f)
			data := []byte(str)
			if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
				http.Error(w, "failed to write file: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
		}

		// finally, write the data stock as response
		stockName := filepath.Base(stock.dir)
		resp := response{
			IsRoot:    stockName == "root",
			ID:        stock.uid,
			ShortName: stockName,
			Name:      stockName,
		}
		writeXML(&resp, w)
	}
}
