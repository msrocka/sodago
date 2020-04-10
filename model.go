package main

import (
	"encoding/xml"
)

const (
	processPath      = "processes"
	flowPath         = "flows"
	flowPropertyPath = "flowproperties"
	unitGroupPath    = "unitgroups"
	contactPath      = "contacts"
	sourcePath       = "sources"
	methodPath       = "lciamethods"
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

// An InfoList contains a list of data set information.
type InfoList struct {
	XMLName        xml.Name           `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetList"`
	TotalSize      int                `xml:"totalSize,attr"`
	StartIndex     int                `xml:"startIndex,attr"`
	PageSize       int                `xml:"pageSize,attr"`
	Processes      []ProcessInfo      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
	Flows          []FlowInfo         `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
	FlowProperties []FlowPropertyInfo `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
	UnitGroups     []UnitGroupInfo    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
	Contacts       []ContactInfo      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
	Sources        []SourceInfo       `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}

// BaseInfo contains some meta data of a data set.
type BaseInfo struct {
	UUID    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI uuid"`
	Version string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetVersion"`
	Name    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI name"`
}

// ProcessInfo contains some meta data of a process data set.
type ProcessInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
}

// FlowInfo contains some meta data of a flow data set.
type FlowInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
}

// ReadFlowInfo reads the flow information from the given flow data set.

// FlowPropertyInfo contains some meta data of a flow property data set.
type FlowPropertyInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
}

// UnitGroupInfo contains some meta data of an unit group data set.
type UnitGroupInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
}

// ContactInfo contains some meta data of a contact data set.
type ContactInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
}

// SourceInfo contains some meta data of a source data set.
type SourceInfo struct {
	BaseInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}
