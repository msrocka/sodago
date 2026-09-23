package main

import (
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

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
		data, err := io.ReadAll(r.Body)
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
			http.Error(w, "unknown data stock", http.StatusBadRequest)
			return
		}
		path := vars["path"]
		if !isValidPath(path) {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		// read the query parameters (all result sets are paged)
		query := r.URL.Query()
		startIndex := queryInt(query, "startIndex", 0)
		if startIndex < 0 {
			startIndex = 0
		}
		pageSize := queryInt(query, "pageSize", defaultPageSize)
		if pageSize <= 0 {
			pageSize = defaultPageSize
		}
		countOnly := queryBool(query, "countOnly")
		allVersions := queryBool(query, "allVersions")
		search := queryBool(query, "search")

		// collect the entries of the requested data set type
		var entries []*indexEntry
		if stock.idx != nil && stock.idx.Entries != nil {
			entries = latestVersions(stock.idx.Entries[path], allVersions)
		}

		// simple search: filter the entries by their name
		if queryBool(query, "search") {
			entries = filterByName(entries, query.Get("name"))
		}
		if search {
			entries = filterByName(entries, query.Get("name"))
		}

		resp := DescriptorList{
			TotalSize:  len(entries),
			StartIndex: startIndex,
			PageSize:   pageSize,
		}
		if countOnly {
			writeXML(&resp, w)
			return
		}

		// select the page entries
		if startIndex > len(entries) {
			startIndex = len(entries)
		}
		end := startIndex + pageSize
		if end > len(entries) {
			end = len(entries)
		}
		for _, e := range entries[startIndex:end] {
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
			case methodPath:
				resp.ImpactCategories = append(resp.ImpactCategories,
					ImpactCategoryDescriptor{BaseDescriptor: base})
			}
		}

		writeXML(&resp, w)
	}
}

// defaultPageSize is the number of data sets that is returned in a list
// response when no page size is specified in the request.
const defaultPageSize = 500

// queryInt returns the integer value of the given query parameter or the given
// default value when the parameter is not set or not a valid integer.
func queryInt(query url.Values, name string, defaultValue int) int {
	value := strings.TrimSpace(query.Get(name))
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

// queryBool returns true when the given query parameter is set to a true value.
func queryBool(query url.Values, name string) bool {
	return strings.EqualFold(strings.TrimSpace(query.Get(name)), "true")
}

func (s *server) handleGetDataSetOverview() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		// select the data stock
		stock := s.dir.findDataStock(vars["datastock"])
		if stock == nil {
			http.Error(w, "unknown data stock", http.StatusBadRequest)
			return
		}

		// select the index entries
		path := vars["path"]
		if !isValidPath(path) {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		entries, ok := stock.idx.Entries[path]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		// find the entry
		uid := vars["id"]
		if uid == "" {
			http.Error(w, "no data set ID given", http.StatusBadRequest)
			return
		}
		var entry *indexEntry
		var version *Version
		for i := range entries {
			e := entries[i]
			if e.UUID != uid {
				continue
			}
			v := ParseVersion(e.Version)
			if entry == nil || v.NewerThan(version) {
				entry = e
				version = v
			}
		}
		if entry == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		base := BaseDescriptor{
			UUID:    uid,
			Name:    entry.Name,
			Version: version.String(),
		}

		switch path {
		case contactPath:
			writeXML(&ContactDescriptor{BaseDescriptor: base}, w)
		case flowPath:
			writeXML(&FlowDescriptor{BaseDescriptor: base}, w)
		case flowPropertyPath:
			writeXML(&FlowPropertyDescriptor{BaseDescriptor: base}, w)
		case methodPath:
			writeXML(&ImpactCategoryDescriptor{BaseDescriptor: base}, w)
		case processPath:
			writeXML(&ProcessDescriptor{BaseDescriptor: base}, w)
		case sourcePath:
			writeXML(&SourceDescriptor{BaseDescriptor: base}, w)
		case unitGroupPath:
			writeXML(&UnitGroupDescriptor{BaseDescriptor: base}, w)

		default:
			http.Error(
				w, "unknown path: "+path, http.StatusInternalServerError)
		}

	}

}
