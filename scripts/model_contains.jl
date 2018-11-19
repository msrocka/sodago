function gencode(data_type, list, doc_type)
  text = """
  // Contains$data_type returns true if there is already a $doc_type
  // with the given ID and version in the list.
  func (list *InfoList) Contains$data_type(other *$(data_type)Info) bool {
  	for i := range list.$list {
  		info := &list.$list[i]
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
  """
end

data = [
  ("Process", "Processes", "process"),
  ("Flow", "Flows", "flow"),
  ("FlowProperty", "FlowProperties", "flow property"),
  ("UnitGroup", "UnitGroups", "unit group"),
  ("Contact", "Contacts", "contact"),
  ("Source", "Sources", "source")
]

for d in data
  println(gencode(d...))
end