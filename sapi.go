package main

import (
	"encoding/xml"
)

// SapiList is the list of data sets returned by the service API.
type SapiList struct {
	XMLName        xml.Name           `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetList"`
	TotalSize      int                `xml:"totalSize,attr"`
	StartIndex     int                `xml:"startIndex,attr"`
	PageSize       int                `xml:"pageSize,attr"`
	Processes      []SapiProcess      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
	Flows          []SapiFlow         `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
	FlowProperties []SapiFlowProperty `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
	UnitGroups     []SapiUnitGroup    `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
	Contacts       []SapiContact      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
	Sources        []SapiSource       `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}

// SapiDataSet contains basic information of each data set type provided by the service API.
type SapiDataSet struct {
	UUID    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI uuid"`
	Version string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetVersion"`
	Name    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI name"`
}

// SapiProcess contains process information provided by the service API.
type SapiProcess struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
}

// SapiFlow contains flow information provided by the service API
type SapiFlow struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
}

// SapiFlowProperty contains flow property information provided by the service API
type SapiFlowProperty struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
}

// SapiUnitGroup contains unit group information provided by the service API
type SapiUnitGroup struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
}

// SapiContact contains contact information provided by the service API
type SapiContact struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
}

// SapiSource contains contact information provided by the service API
type SapiSource struct {
	SapiDataSet
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}
