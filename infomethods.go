package main

import (
	"net/http"
)

// DataSetInfo contains some common methods for the info data of different
// data set types
type DataSetInfo interface {
	// Key returns the key under which the data set ist stored.
	Key(stock *DataStock) string

	// Base returns the base information of the data set.
	Base() *BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *ProcessInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Process/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *ProcessInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *FlowInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Flow/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *FlowInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *FlowPropertyInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/FlowProperty/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *FlowPropertyInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *UnitGroupInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/UnitGroup/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *UnitGroupInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *ContactInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Contact/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *ContactInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Key returns the key under which the data set ist stored.
func (info *SourceInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Source/" + info.UUID + "/" + v
}

// Base returns the base information of the data set.
func (info *SourceInfo) Base() *BaseInfo {
	return &info.BaseInfo
}

// Matches returns true if the given info object matches the given ID
// and version. The given version may be nil which means that it is
// ignored in the check.
func Matches(info DataSetInfo, id string, version *Version) bool {
	if info == nil || info.Base().UUID != id {
		return false
	}
	if version == nil {
		return true
	}
	v := ParseVersion(info.Base().Version)
	return version.Same(v)
}

// ServeDataSet loads the data set described by the given info object
// from the given data stock  and writes it to the response.
func ServeDataSet(info DataSetInfo, stock *DataStock,
	w http.ResponseWriter) {
	if info == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	data := db.Get(DataSetBucket, info.Key(stock))
	if data == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	ServeXMLBytes(data, w)
}
