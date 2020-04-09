package main

import (
	"encoding/xml"
	"net/http"
	"path/filepath"
)

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
		ServeXML(&resp, w)
	}
}
