package main

// ContainsProcess returns true if there is already a process
// with the given ID and version in the list.
func (list *InfoList) ContainsProcess(other *ProcessInfo) bool {
	for i := range list.Processes {
		info := &list.Processes[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}

// ContainsFlow returns true if there is already a flow
// with the given ID and version in the list.
func (list *InfoList) ContainsFlow(other *FlowInfo) bool {
	for i := range list.Flows {
		info := &list.Flows[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}

// ContainsFlowProperty returns true if there is already a flow property
// with the given ID and version in the list.
func (list *InfoList) ContainsFlowProperty(other *FlowPropertyInfo) bool {
	for i := range list.FlowProperties {
		info := &list.FlowProperties[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}

// ContainsUnitGroup returns true if there is already a unit group
// with the given ID and version in the list.
func (list *InfoList) ContainsUnitGroup(other *UnitGroupInfo) bool {
	for i := range list.UnitGroups {
		info := &list.UnitGroups[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}

// ContainsContact returns true if there is already a contact
// with the given ID and version in the list.
func (list *InfoList) ContainsContact(other *ContactInfo) bool {
	for i := range list.Contacts {
		info := &list.Contacts[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}

// ContainsSource returns true if there is already a source
// with the given ID and version in the list.
func (list *InfoList) ContainsSource(other *SourceInfo) bool {
	for i := range list.Sources {
		info := &list.Sources[i]
		if info.UUID != other.UUID {
			continue
		}
		v1 := ParseVersion(info.Version)
		v2 := ParseVersion(other.Version)
		if v1.Compare(v2) == 0 {
			return true
		}
	}
	return false
}
