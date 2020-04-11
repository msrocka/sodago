package main

import (
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (s *server) registerRoutes(r *mux.Router, args *Args) {

	// data stocks
	r.HandleFunc("/resource/datastocks",
		s.handleGetDataStocks()).Methods("GET")

	// profiles
	r.HandleFunc("/resource/profiles/{id}",
		s.handleGetProfile()).Methods("GET")
	r.HandleFunc("/resource/profiles",
		s.handleGetProfiles()).Methods("GET")

	// login: currently nothing is checked here
	r.HandleFunc("/resource/authenticate/login",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Logged in"))
		})

	// GET a single data set from the root data stock
	// specific version
	r.HandleFunc("/resource/{path}/{id}",
		s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/{path}/{id}",
		s.handleGetDataSet()).
		Methods("GET", "HEAD")

	// GET a single data set from a data stock
	// specific version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).Methods("GET", "HEAD")

	// POST a data set
	r.HandleFunc("/resource/{path}",
		s.handlePostDataSet()).Methods("POST")

	// GET  a list of data sets
	r.HandleFunc("/resource/datastocks/{datastock}/{path}",
		s.handleGetDataSets()).Methods("GET")
	r.HandleFunc("/resource/{path}",
		s.handleGetDataSets()).Methods("GET")
}

func (s *server) handleGetDataStocks() http.HandlerFunc {

	type item struct {
		IsRoot    bool   `xml:"root,attr"`
		ID        string `xml:"uuid"`
		ShortName string `xml:"shortName"`
		Name      string `xml:"name"`
	}

	type response struct {
		XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataStockList"`
		Items   []item   `xml:"dataStock"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{}
		for _, stock := range s.dir.dataStocks {
			name := filepath.Base(stock.dir)
			resp.Items = append(resp.Items, item{
				IsRoot:    name == "root",
				ID:        stock.uid,
				ShortName: name,
				Name:      name,
			})
		}
		writeXML(&resp, w)
	}
}
