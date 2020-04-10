package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
)

// An index stores the basic data set informations in a map
// path -> entries
type index struct {
	Entries map[string][]*indexEntry `json:"entries"`
}

type indexEntry struct {
	UUID    string `json:"uuid"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

func (idx *index) save(file string) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, os.ModePerm)
}

func readIndex(file string) (*index, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	idx := &index{}
	if err := json.Unmarshal(data, idx); err != nil {
		return nil, err
	}
	return idx, nil
}

func (idx *index) contains(path string, entry *indexEntry) bool {
	for _, e := range idx.Entries[path] {
		if e.UUID == entry.UUID && e.Version == entry.Version {
			return true
		}
	}
	return false
}

// Reads the index information from the raw XML bytes of the given
// data set. The path is the request path for the respective data
// set type.
func extractIndexEntry(path string, dataSet []byte) (*indexEntry, error) {

	if path == processPath {
		d := &struct {
			XMLName xml.Name `xml:"processDataSet"`
			Name    string   `xml:"processInformation>dataSetInformation>name>baseName"`
			UUID    string   `xml:"processInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	if path == flowPath {
		d := &struct {
			XMLName xml.Name `xml:"flowDataSet"`
			Name    string   `xml:"flowInformation>dataSetInformation>name>baseName"`
			UUID    string   `xml:"flowInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	if path == flowPropertyPath {
		d := &struct {
			XMLName xml.Name `xml:"flowPropertyDataSet"`
			Name    string   `xml:"flowPropertiesInformation>dataSetInformation>name"`
			UUID    string   `xml:"flowPropertiesInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	if path == unitGroupPath {
		d := &struct {
			XMLName xml.Name `xml:"unitGroupDataSet"`
			Name    string   `xml:"unitGroupInformation>dataSetInformation>name"`
			UUID    string   `xml:"unitGroupInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	if path == contactPath {
		d := &struct {
			XMLName xml.Name `xml:"contactDataSet"`
			Name    string   `xml:"contactInformation>dataSetInformation>name"`
			UUID    string   `xml:"contactInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	if path == sourcePath {
		d := &struct {
			XMLName xml.Name `xml:"sourceDataSet"`
			Name    string   `xml:"sourceInformation>dataSetInformation>shortName"`
			UUID    string   `xml:"sourceInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: d.Name, Version: d.Version}, nil
	}

	// TODO: LCIA methods / categories

	return nil, errors.New("unknown path: " + path)
}
