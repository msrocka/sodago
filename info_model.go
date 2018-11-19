package main

import (
	"encoding/xml"
	"log"
)

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

// DataSetInfo contains some meta data of a data set.
type DataSetInfo struct {
	UUID    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI uuid"`
	Version string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetVersion"`
	Name    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI name"`
}

// ProcessInfo contains some meta data of a process data set.
type ProcessInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
}

// FlowInfo contains some meta data of a flow data set.
type FlowInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
}

func ReadFlowInfo(data []byte) *FlowInfo {
	info := &struct {
		XMLName xml.Name `xml:"flowDataSet"`
		Name    string   `xml:"flowInformation>dataSetInformation>name>baseName"`
		UUID    string   `xml:"flowInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, info)
	if err != nil {
		log.Println("ERROR: failed to read flow info", err)
		return nil
	}
	flowInfo := &FlowInfo{}
	flowInfo.Name = info.Name
	flowInfo.UUID = info.UUID
	flowInfo.Version = info.Version
	return flowInfo
}

// FlowPropertyInfo contains some meta data of a flow property data set.
type FlowPropertyInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
}

// UnitGroupInfo contains some meta data of an unit group data set.
type UnitGroupInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
}

// ContactInfo contains some meta data of a contact data set.
type ContactInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
}

// SourceInfo contains some meta data of a source data set.
type SourceInfo struct {
	DataSetInfo
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}
