package main

import (
	"log"
	"os"
	"path/filepath"
)

type datadir struct {
	root string
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

	dir := &datadir{root}

	// create the root data stock if necessary
	if !dir.isDataStock("root") {
		if err := dir.createDataStock("root"); err != nil {
			return nil, err
		}
	}

	return dir, nil
}

func (dir *datadir) createDataStock(name string) error {
	log.Println("Create data stock", name)

	return nil
}

func (dir *datadir) isDataStock(name string) bool {
	path := filepath.Join(dir.root, name)
	if !fileExists(path) {
		return false
	}
	meta := filepath.Join(path, "meta.xml")
	if !fileExists(meta) {
		return false
	}
	return true
}
