package main

// Key returns the key under which the data set ist stored.
func (info *ProcessInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Process/" + info.UUID + "/" + v
}

// Key returns the key under which the data set ist stored.
func (info *FlowInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Flow/" + info.UUID + "/" + v
}

// Key returns the key under which the data set ist stored.
func (info *FlowPropertyInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/FlowProperty/" + info.UUID + "/" + v
}

// Key returns the key under which the data set ist stored.
func (info *UnitGroupInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/UnitGroup/" + info.UUID + "/" + v
}

// Key returns the key under which the data set ist stored.
func (info *ContactInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Contact/" + info.UUID + "/" + v
}

// Key returns the key under which the data set ist stored.
func (info *SourceInfo) Key(stock *DataStock) string {
	v := ParseVersion(info.Version).String() // standard format
	return stock.ID + "/Source/" + info.UUID + "/" + v
}
