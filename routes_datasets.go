package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (s *server) handlePostDataSet() http.HandlerFunc {

	type response struct {
		IsRoot    bool   `xml:"root,attr"`
		ID        string `xml:"uuid"`
		ShortName string `xml:"shortName"`
		Name      string `xml:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read body "+err.Error(), http.StatusBadRequest)
			return
		}
		r.Body.Close()
		stockID := r.Header.Get("stock")
		path := mux.Vars(r)["path"]

		// save the data set
		s.mutex.Lock()
		stock, err := s.dir.put(stockID, path, data)
		s.mutex.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// return the data stock on success
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

func (s *server) handleGetDataSet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		stockID := vars["datastock"]
		path := vars["path"]
		uid := vars["id"]
		version := vars["version"]

		data, err := s.dir.get(stockID, path, &indexEntry{
			UUID:    uid,
			Version: version,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeBytesXML(data, w)
	}
}

func (s *server) handleGetDataSets() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// check data stock and path
		vars := mux.Vars(r)
		stock := s.dir.findDataStock(vars["datastock"])
		if stock == nil {
			http.Error(w, "Unknown data stock", http.StatusBadRequest)
			return
		}
		path := vars["path"]
		if !isValidPath(path) {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		resp := DescriptorList{}
		if stock.idx == nil || stock.idx.Entries == nil {
			writeXML(&resp, w)
			return
		}
		entries, ok := stock.idx.Entries[path]
		if !ok {
			writeXML(&resp, w)
			return
		}

		resp.PageSize = len(entries)
		resp.TotalSize = len(entries)
		resp.StartIndex = 0

		for _, e := range entries {
			base := BaseDescriptor{
				UUID:    e.UUID,
				Name:    e.Name,
				Version: e.Version,
			}
			switch path {
			case processPath:
				resp.Processes = append(resp.Processes, ProcessDescriptor{BaseDescriptor: base})
			case flowPath:
				resp.Flows = append(resp.Flows, FlowDescriptor{BaseDescriptor: base})
			case flowPropertyPath:
				resp.FlowProps = append(resp.FlowProps, FlowPropertyDescriptor{BaseDescriptor: base})
			case unitGroupPath:
				resp.UnitGroups = append(resp.UnitGroups, UnitGroupDescriptor{BaseDescriptor: base})
			case contactPath:
				resp.Contacts = append(resp.Contacts, ContactDescriptor{BaseDescriptor: base})
			case sourcePath:
				resp.Sources = append(resp.Sources, SourceDescriptor{BaseDescriptor: base})
			}
		}

		writeXML(&resp, w)
	}
}
