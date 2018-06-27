package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

// DataStockList contains a list of data stocks. This type is used for XML
// serialization.
type DataStockList struct {
	XMLName    xml.Name    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataStockList"`
	DataStocks []DataStock `xml:"dataStock"`
}

// A DataStock contains a set if data sets.
type DataStock struct {
	IsRoot      bool   `xml:"root,attr"`
	ID          string `xml:"uuid"`
	ShortName   string `xml:"shortName"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

func (stock *DataStock) String() string {
	return stock.Name + "@" + stock.ID
}

// InitStocks loads the data stocks from the data folder. It creates a root
// data stock if it does not exist yet.
func InitStocks() []*DataStock {
	stocks := GetStockInfos()
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

// GetStockInfos returns the meta data from the data stocks.
func GetStockInfos() []*DataStock {
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
		stock := GetStockInfo("data/stocks/" + subDir.Name())
		if stock != nil {
			stocks = append(stocks, stock)
		}
	}
	return stocks
}

// GetStockInfo reads the data stock information from the 'meta.xml'
// file of the given folder. It returns nil if there is no meta.xml file
// or if reading this file failed.
func GetStockInfo(folder string) *DataStock {
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

// writeStockInfo writes the given data stock information to the meta.xml
// file of the given folder.
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
