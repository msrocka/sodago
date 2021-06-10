package main

import "encoding/xml"

type BaseDescriptor struct {
	UUID    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI uuid"`
	Version string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetVersion"`
	Name    string `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI name"`
}

type ProcessDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
}

type FlowDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
}

type FlowPropertyDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
}

type UnitGroupDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
}

type ContactDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
}

type SourceDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
}

type ImpactCategoryDescriptor struct {
	BaseDescriptor
	XMLName xml.Name `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/LCIAMethod LCIAMethod"`
}

type DescriptorList struct {
	XMLName          xml.Name                   `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI dataSetList"`
	TotalSize        int                        `xml:"totalSize,attr"`
	StartIndex       int                        `xml:"startIndex,attr"`
	PageSize         int                        `xml:"pageSize,attr"`
	Processes        []ProcessDescriptor        `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Process process"`
	Flows            []FlowDescriptor           `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Flow flow"`
	FlowProps        []FlowPropertyDescriptor   `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/FlowProperty flowProperty"`
	UnitGroups       []UnitGroupDescriptor      `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/UnitGroup unitGroup"`
	Contacts         []ContactDescriptor        `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Contact contact"`
	Sources          []SourceDescriptor         `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/Source source"`
	ImpactCategories []ImpactCategoryDescriptor `xml:"http://www.ilcd-network.org/ILCD/ServiceAPI/LCIAMethod LCIAMethod"`
}
