package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

// DataStockList contains a list of data stocks. This type is used for XML
// serialization.
type DataStockList struct {
	XMLName    xml.Name     `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataStockList"`
	DataStocks []*DataStock `xml:"dataStock"`
}

// A DataStock contains a set if data sets.
type DataStock struct {
	IsRoot      bool   `xml:"root,attr"`
	ID          string `xml:"uuid"`
	ShortName   string `xml:"shortName"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

// GetDataStocksHandler returns a handler function for returning the data stock
// meta data.
func GetDataStocksHandler(context *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stockList := DataStockList{DataStocks: context.DataStocks}
		ServeXML(&stockList, w)
	}
}

func (stock *DataStock) String() string {
	return stock.Name + "@" + stock.ID
}

// InitStocks loads the data stocks from the data folder. It creates a root
// data stock if it does not exist yet.
func InitStocks() []*DataStock {
	stocks := readStockInfos()
	hasRoot := false
	for _, info := range stocks {
		if info.IsRoot {
			hasRoot = true
			break
		}
	}
	if !hasRoot {
		root := &DataStock{
			IsRoot:      true,
			ID:          uuid.NewV4().String(),
			ShortName:   "root",
			Name:        "root",
			Description: "The root data stock"}
		os.MkdirAll("data/stocks/root", os.ModePerm)
		writeStockInfo(root, "data/stocks/root")
		stocks = append(stocks, root)
	}
	return stocks
}

func readStockInfos() []*DataStock {
	var stocks []*DataStock
	stat, err := os.Stat("data/stocks")
	if err != nil || !stat.IsDir() {
		return stocks
	}
	subDirs, err := ioutil.ReadDir("data/stocks")
	if err != nil {
		log.Println("ERROR: failed to read data stock folders", err.Error())
		return stocks
	}
	for _, subDir := range subDirs {
		stock := readStockInfo("data/stocks/" + subDir.Name())
		if stock != nil {
			stocks = append(stocks, stock)
		}
	}
	return stocks
}

func readStockInfo(folder string) *DataStock {
	path := filepath.Join(folder, "meta.xml")
	if _, err := os.Stat(path); err != nil {
		log.Println("No data stock info found at", path)
		return nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("ERROR: Failed to read file", path)
		return nil
	}
	stock := &DataStock{}
	err = xml.Unmarshal(data, stock)
	if err != nil {
		log.Println("ERROR: Failed to read file", path)
		return nil
	}
	return stock
}

func writeStockInfo(stock *DataStock, folder string) {
	data, err := xml.MarshalIndent(&stock, "", "    ")
	if err != nil {
		log.Println("ERROR: failed to marshal data stock", stock)
		return
	}
	path := filepath.Join(folder, "meta.xml")
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		log.Println("ERROR: failed to write file", path)
	}
}
