package main

import (
	"encoding/xml"
	"net/http"
	"path/filepath"
)

func (s *server) handleGetDataStocks() http.HandlerFunc {

	type StockItem struct {
		IsRoot    bool   `xml:"root,attr"`
		ID        string `xml:"uuid"`
		ShortName string `xml:"shortName"`
		Name      string `xml:"name"`
	}

	type Response struct {
		XMLName xml.Name    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataStockList"`
		Items   []StockItem `xml:"dataStock"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp := Response{}
		for _, stock := range s.dir.dataStocks {
			name := filepath.Base(stock.dir)
			resp.Items = append(resp.Items, StockItem{
				IsRoot:    name == "root",
				ID:        stock.uid,
				ShortName: name,
				Name:      name,
			})
		}
		ServeXML(&resp, w)
	}
}
