package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
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

// firstName returns the first entry of a list of names. Names of an ILCD data
// set are given in multiple languages where the first name is the default one.
func firstName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
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
			Names   []string `xml:"processInformation>dataSetInformation>name>baseName"`
			UUID    string   `xml:"processInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: firstName(d.Names), Version: d.Version}, nil
	}

	if path == flowPath {
		d := &struct {
			XMLName xml.Name `xml:"flowDataSet"`
			Names   []string `xml:"flowInformation>dataSetInformation>name>baseName"`
			UUID    string   `xml:"flowInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: firstName(d.Names), Version: d.Version}, nil
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
		// an LCIA method data set stores the name directly in the `name`
		// element of the data set information (see the ILCD LCIA method
		// schema)
		d := &struct {
			XMLName xml.Name `xml:"LCIAMethodDataSet"`
			Names   []string `xml:"LCIAMethodInformation>dataSetInformation>name"`
			UUID    string   `xml:"LCIAMethodInformation>dataSetInformation>UUID"`
			Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
		}{}
		if err := xml.Unmarshal(dataSet, d); err != nil {
			return nil, err
		}
		return &indexEntry{UUID: d.UUID, Name: firstName(d.Names), Version: d.Version}, nil
	}

	return nil, errors.New("unknown path: " + path)
}
