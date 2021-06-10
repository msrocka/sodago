package main

import (
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func (s *server) registerRoutes(r *mux.Router) {

	// GET data stocks
	r.HandleFunc("/resource/datastocks",
		s.handleGetDataStocks()).Methods("GET")

	// profiles
	// GET profile
	r.HandleFunc("/resource/profiles/{id}",
		s.handleGetProfile()).Methods("GET")
	// GET profiles
	r.HandleFunc("/resource/profiles",
		s.handleGetProfiles()).Methods("GET")

	// authentication
	r.HandleFunc("/resource/authenticate/login", s.handleGetLogin()).
		Queries("userName", "{user}", "password", "{password}")
	r.HandleFunc("/resource/authenticate/status",
		s.handleGetAuthenticationStatus())

	// GET a single data set from the root data stock
	// get an overview
	r.HandleFunc("/resource/{path}/{id}", s.handleGetDataSetOverview()).
		Queries("view", "overview").
		Methods("GET", "HEAD")
	// specific version
	r.HandleFunc("/resource/{path}/{id}", s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/{path}/{id}", s.handleGetDataSet()).
		Methods("GET", "HEAD")

	// GET a digital/external file of a source
	r.HandleFunc("/resource/datastocks/{datastock}/sources/{id}/{file}",
		s.handleGetExternalFile()).Methods("GET", "HEAD")
	r.HandleFunc("/resource/sources/{id}/{file}",
		s.handleGetExternalFile()).Methods("GET", "HEAD")

	// GET a single data set from a data stock
	// get an overview
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSetOverview()).
		Queries("view", "overview").
		Methods("GET", "HEAD")
	// specific version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).
		Queries("version", "{version}").
		Methods("GET", "HEAD")
	// the latest version
	r.HandleFunc("/resource/datastocks/{datastock}/{path}/{id}",
		s.handleGetDataSet()).Methods("GET", "HEAD")

	// POST a data set
	r.HandleFunc("/resource/sources/withBinaries",
		s.handlePostSourceWithFiles()).Methods("POST")
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
