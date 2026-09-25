package main

import (
	"encoding/xml"
	"mime"
	"net/http"
	"net/url"
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

		// the `digitalfile` route returns the first digital file of a source
		if file == "digitalfile" {
			file = s.firstDigitalFile(stock, uid)
			if file == "" {
				http.Error(w, "Unknown file", http.StatusNotFound)
				return
			}
		}

		path := filepath.Join(stock.dir, "external_docs", uid, file)
		if !fileExists(path) {
			http.Error(w, "Unknown file", http.StatusNotFound)
			return
		}

		data, err := os.ReadFile(path)
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

// firstDigitalFile returns the name of the first digital file that is
// referenced in the given source data set or an empty string when the source
// has no such reference.
func (s *server) firstDigitalFile(stock *dataStock, uid string) string {
	data, err := s.dir.get(stock.uid, sourcePath, &indexEntry{UUID: uid})
	if err != nil {
		return ""
	}
	return firstDigitalFileName(data)
}

// firstDigitalFileName reads the file name of the first digital file reference
// from the given source data set XML.
func firstDigitalFileName(source []byte) string {
	d := &struct {
		Files []struct {
			URI string `xml:"uri,attr"`
		} `xml:"sourceInformation>dataSetInformation>referenceToDigitalFile"`
	}{}
	if err := xml.Unmarshal(source, d); err != nil {
		return ""
	}
	for _, file := range d.Files {
		if name := fileNameOf(file.URI); name != "" {
			return name
		}
	}
	return ""
}

// fileNameOf returns the file name of a digital file reference. The last part
// of the URI is used and unescaped (like the clients do it).
func fileNameOf(uri string) string {
	name := strings.TrimSpace(strings.ReplaceAll(uri, `\`, "/"))
	if pos := strings.LastIndex(name, "/"); pos >= 0 {
		name = name[pos+1:]
	}
	if name == "" {
		return ""
	}
	if decoded, err := url.QueryUnescape(name); err == nil {
		return decoded
	}
	return name
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
			if err := os.WriteFile(path, data, os.ModePerm); err != nil {
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
