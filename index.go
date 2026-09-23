package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strings"
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
	return os.WriteFile(file, data, os.ModePerm)
}

func readIndex(file string) (*index, error) {
	data, err := os.ReadFile(file)
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

// latestVersions returns only the most recent version of each data set when
// allVersions is false, which is the default of a list request. The order of
// the entries is kept.
func latestVersions(entries []*indexEntry, allVersions bool) []*indexEntry {
	if allVersions {
		return entries
	}
	latest := make(map[string]*indexEntry)
	order := make([]string, 0, len(entries))
	for _, e := range entries {
		current, ok := latest[e.UUID]
		if !ok {
			latest[e.UUID] = e
			order = append(order, e.UUID)
			continue
		}
		if ParseVersion(e.Version).NewerThan(ParseVersion(current.Version)) {
			latest[e.UUID] = e
		}
	}
	result := make([]*indexEntry, 0, len(order))
	for _, uid := range order {
		result = append(result, latest[uid])
	}
	return result
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

	if path == methodPath {
		return extractMethodEntry(dataSet)
	}

	return nil, errors.New("unknown path: " + path)
}

// Reads the index information from the raw XML bytes of an LCIA method data
// set. Depending on the data source, the element that contains the data set
// information is written as `LCIAMethodInformation` or `lciaMethodInformation`,
// thus the UUID, name and version are read via the local element names.
func extractMethodEntry(dataSet []byte) (*indexEntry, error) {
	decoder := xml.NewDecoder(bytes.NewReader(dataSet))
	entry := &indexEntry{}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		var target *string
		switch start.Name.Local {
		case "UUID":
			target = &entry.UUID
		case "baseName":
			target = &entry.Name
		case "dataSetVersion":
			target = &entry.Version
		default:
			continue
		}
		if *target != "" {
			continue
		}
		value := ""
		if err := decoder.DecodeElement(&value, &start); err != nil {
			return nil, err
		}
		*target = strings.TrimSpace(value)
	}
	return entry, nil
}
