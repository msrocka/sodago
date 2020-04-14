package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
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
		stock, err := s.dir.put(stockID, path, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

func (s *server) handlePostSourceWithFiles() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "failed to parse multi-part form: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// try to read the source
		source, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "failed to get `file` param: "+err.Error(),
				http.StatusBadRequest)
			return
		}
		defer source.Close()
		data, err := ioutil.ReadAll(source)
		if err != nil {
			http.Error(w, "failed to read `file` param: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// extract the source info
		sourceInfo, err := extractIndexEntry(sourcePath, data)
		if err != nil {
			http.Error(w, "failed to read `file` param: "+err.Error(),
				http.StatusBadRequest)
			return
		}

		// try to save the source
		stockID := r.Header.Get("stock")
		stock, err := s.dir.put(stockID, sourcePath, data)
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
		for f := range r.MultipartForm.File {
			if f == "file" {
				continue
			}
			file, _, err := r.FormFile(f)
			if err != nil {
				http.Error(w, "failed to get file param: "+err.Error(),
					http.StatusBadRequest)
				return
			}
			defer file.Close()
			data, err := ioutil.ReadAll(file)
			if err != nil {
				http.Error(w, "failed to read file param: "+err.Error(),
					http.StatusBadRequest)
				return
			}
			path := filepath.Join(dir, f)
			if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
				http.Error(w, "failed to write file: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
		}
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

	type base struct {
		UUID    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI uuid"`
		Version string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetVersion"`
		Name    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI name"`
	}

	type process struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
	}

	type flow struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
	}

	type flowProp struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
	}

	type unitGroup struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
	}

	type contact struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
	}

	type source struct {
		base
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
	}

	type response struct {
		XMLName    xml.Name    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetList"`
		TotalSize  int         `xml:"totalSize,attr"`
		StartIndex int         `xml:"startIndex,attr"`
		PageSize   int         `xml:"pageSize,attr"`
		Processes  []process   `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
		Flows      []flow      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
		FlowProps  []flowProp  `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
		UnitGroups []unitGroup `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
		Contacts   []contact   `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
		Sources    []source    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
	}

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

		resp := response{}
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
			base := base{
				UUID:    e.UUID,
				Name:    e.Name,
				Version: e.Version,
			}
			switch path {
			case processPath:
				resp.Processes = append(resp.Processes, process{base: base})
			case flowPath:
				resp.Flows = append(resp.Flows, flow{base: base})
			case flowPropertyPath:
				resp.FlowProps = append(resp.FlowProps, flowProp{base: base})
			case unitGroupPath:
				resp.UnitGroups = append(resp.UnitGroups, unitGroup{base: base})
			case contactPath:
				resp.Contacts = append(resp.Contacts, contact{base: base})
			case sourcePath:
				resp.Sources = append(resp.Sources, source{base: base})
			}
		}

		writeXML(&resp, w)
	}
}
