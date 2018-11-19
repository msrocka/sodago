package main

import (
	"encoding/xml"
	"log"
)

// ReadProcessInfo reads the meta data of the given process data set.
func ReadProcessInfo(data []byte) *ProcessInfo {
	d := &struct {
		XMLName xml.Name `xml:"processDataSet"`
		Name    string   `xml:"processInformation>dataSetInformation>name>baseName"`
		UUID    string   `xml:"processInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read ProcessInfo", err)
		return nil
	}
	info := &ProcessInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}

// ReadFlowInfo reads the meta data of the given flow data set.
func ReadFlowInfo(data []byte) *FlowInfo {
	d := &struct {
		XMLName xml.Name `xml:"flowDataSet"`
		Name    string   `xml:"flowInformation>dataSetInformation>name>baseName"`
		UUID    string   `xml:"flowInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read FlowInfo", err)
		return nil
	}
	info := &FlowInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}

// ReadFlowPropertyInfo reads the meta data of the given flow property data set.
func ReadFlowPropertyInfo(data []byte) *FlowPropertyInfo {
	d := &struct {
		XMLName xml.Name `xml:"flowPropertyDataSet"`
		Name    string   `xml:"flowPropertiesInformation>dataSetInformation>name"`
		UUID    string   `xml:"flowPropertiesInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read FlowPropertyInfo", err)
		return nil
	}
	info := &FlowPropertyInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}

// ReadUnitGroupInfo reads the meta data of the given unit group data set.
func ReadUnitGroupInfo(data []byte) *UnitGroupInfo {
	d := &struct {
		XMLName xml.Name `xml:"unitGroupDataSet"`
		Name    string   `xml:"unitGroupInformation>dataSetInformation>name"`
		UUID    string   `xml:"unitGroupInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read UnitGroupInfo", err)
		return nil
	}
	info := &UnitGroupInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}

// ReadContactInfo reads the meta data of the given contact data set.
func ReadContactInfo(data []byte) *ContactInfo {
	d := &struct {
		XMLName xml.Name `xml:"contactDataSet"`
		Name    string   `xml:"contactInformation>dataSetInformation>name"`
		UUID    string   `xml:"contactInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read ContactInfo", err)
		return nil
	}
	info := &ContactInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}

// ReadSourceInfo reads the meta data of the given source data set.
func ReadSourceInfo(data []byte) *SourceInfo {
	d := &struct {
		XMLName xml.Name `xml:"sourceDataSet"`
		Name    string   `xml:"sourceInformation>dataSetInformation>shortName"`
		UUID    string   `xml:"sourceInformation>dataSetInformation>UUID"`
		Version string   `xml:"administrativeInformation>publicationAndOwnership>dataSetVersion"`
	}{}
	err := xml.Unmarshal(data, d)
	if err != nil {
		log.Println("ERROR: failed to read SourceInfo", err)
		return nil
	}
	info := &SourceInfo{}
	info.Name = d.Name
	info.UUID = d.UUID
	info.Version = d.Version
	return info
}
