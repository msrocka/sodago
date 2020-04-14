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
	idx *index
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

		// read the data stock ID
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

		// read the index
		idxFile := filepath.Join(stock.dir, "index.json")
		if fileExists(idxFile) {
			idx, err := readIndex(idxFile)
			if err != nil {
				return nil, err
			}
			stock.idx = idx
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

// Get the data stock with the given UUID. The UUID may be empty and we return
// the root data stock in this case.
func (dir *datadir) findDataStock(uid string) *dataStock {
	for i := range dir.dataStocks {
		stock := &dir.dataStocks[i]
		if stock.uid == uid {
			return stock
		}
		if uid == "" {
			name := filepath.Base(stock.dir)
			if name == "root" {
				return stock
			}
		}
	}
	return nil
}

func (dir *datadir) put(stockID string, path string, dataSet []byte) (*dataStock, error) {

	// check data stock and path
	stock := dir.findDataStock(stockID)
	if stock == nil {
		return nil, errUnknownDataStock
	}
	if !isValidPath(path) {
		return nil, errInvalidPath
	}

	// check the index entry
	entry, err := extractIndexEntry(path, dataSet)
	if err != nil {
		return nil, errInvalidDataSet
	}
	if entry.UUID == "" || entry.Version == "" {
		return nil, errInvalidDataSet
	}
	if stock.idx != nil && stock.idx.contains(path, entry) {
		return nil, errDataSetExists
	}

	// store the file
	fileDir := filepath.Join(stock.dir, path)
	if !fileExists(fileDir) {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	file := filepath.Join(fileDir, entry.UUID+"_"+entry.Version+".xml")
	if err := ioutil.WriteFile(file, dataSet, os.ModePerm); err != nil {
		return nil, err
	}

	// register the index entry
	if stock.idx == nil {
		stock.idx = &index{}
	}
	if stock.idx.Entries == nil {
		stock.idx.Entries = make(map[string][]*indexEntry)
	}
	stock.idx.Entries[path] = append(stock.idx.Entries[path], entry)
	idxFile := filepath.Join(stock.dir, "index.json")
	if err := stock.idx.save(idxFile); err != nil {
		return nil, err
	}

	return stock, nil
}

func (dir *datadir) get(stockID string, path string, entry *indexEntry) ([]byte, error) {

	// check data stock and path
	stock := dir.findDataStock(stockID)
	if stock == nil {
		return nil, errUnknownDataStock
	}
	if !isValidPath(path) {
		return nil, errInvalidPath
	}

	// find the file
	fileDir := filepath.Join(stock.dir, path)
	if !fileExists(fileDir) {
		return nil, errDataSetNotExists
	}
	file := ""
	if entry.Version != "" {
		file = entry.UUID + "_" + entry.Version + ".xml"
	} else {

		// if no version is given, we want the latest one
		var version *Version
		prefix := entry.UUID + "_"
		files, err := ioutil.ReadDir(fileDir)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if !strings.HasPrefix(f.Name(), prefix) {
				continue
			}
			v := ParseVersion(strings.TrimPrefix(
				strings.TrimSuffix(f.Name(), ".xml"), prefix))
			if file == "" ||
				version == nil ||
				v.NewerThan(version) {
				file = f.Name()
				version = v
				continue
			}
		}
	}

	// read the file if it exists
	if file == "" {
		return nil, errDataSetNotExists
	}
	file = filepath.Join(fileDir, file)
	if !fileExists(file) {
		return nil, errDataSetNotExists
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return data, err
}
