package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type datadir struct {
	root       string
	dataStocks []dataStock
}

type dataStock struct {
	dir string
	uid string
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func newDataDir(root string) (*datadir, error) {

	// initialize the root folder
	if !fileExists(root) {
		log.Println("initialize data folder @", root)
		if err := os.MkdirAll(root, os.ModePerm); err != nil {
			return nil, err
		}
	}

	dir := &datadir{root: root}

	// read the data stocks
	var rootStock *dataStock
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		meta := filepath.Join(root, file.Name(), ".stock")
		if !fileExists(meta) {
			continue
		}
		data, err := ioutil.ReadFile(meta)
		if err != nil {
			return nil, err
		}
		uid := strings.TrimSpace(string(data))
		stock := dataStock{
			dir: filepath.Join(root, file.Name()),
			uid: uid,
		}
		dir.dataStocks = append(dir.dataStocks, stock)
		if file.Name() == "root" {
			rootStock = &stock
		}
	}

	// create the root data stock if necessary
	if rootStock == nil {
		rootStock, err := dir.createDataStock("root")
		if err != nil {
			return nil, err
		}
		dir.dataStocks = append(dir.dataStocks, *rootStock)
	}

	return dir, nil
}

func (dir *datadir) createDataStock(name string) (*dataStock, error) {
	log.Println("Create data stock", name)
	stockDir := filepath.Join(dir.root, name)
	if err := os.MkdirAll(stockDir, os.ModePerm); err != nil {
		return nil, err
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	meta := filepath.Join(stockDir, ".stock")
	uidStr := uid.String()
	err = ioutil.WriteFile(meta, []byte(uidStr), os.ModePerm)
	if err != nil {
		return nil, err
	}
	stock := dataStock{
		uid: uidStr,
		dir: stockDir,
	}

	return &stock, nil
}
