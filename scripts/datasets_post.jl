function gencode(data_type, list)
  """
  func post$data_type(data []byte, stock *DataStock, w http.ResponseWriter) {
    info := Read$(data_type)Info(data)
    if info == nil {
    	http.Error(w, "Could not read body", http.StatusBadRequest)
    	return
    }
    content := db.Content(stock)
    if content.Contains$(data_type)(info) {
    	http.Error(w, "$(data_type) "+info.UUID+" "+info.Version+
    			"already exists", http.StatusBadRequest)
    	return
    }
    v := ParseVersion(info.Version).String() // standard format
    key := stock.ID + "/$(data_type)/" + info.UUID + "/" + v
    db.Put(DataSetBucket, key, data)
    content.$list = append(content.$list, *info)
    db.UpdateContent(stock, content)
    ServeXML(stock, w)
  }
  """
end

data = [
  ("Process", "Processes"),
  ("Flow", "Flows"),
  ("FlowProperty", "FlowProperties"),
  ("UnitGroup", "UnitGroups"),
  ("Contact", "Contacts"),
  ("Source", "Sources")
]

for d in data
  println(gencode(d...))
end

for d in data
  t = """
  case "$(lowercase(d[2]))":
    post$(d[1])(data, stock, w)
  """
  println(t)
end