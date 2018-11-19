
function gencode(info_type, doc_type)
  text = """
  // Read$info_type reads the meta data of the given $doc_type data set.
  func Read$info_type(data []byte) *$info_type {
    d := &struct {
      XMLName xml.Name `xml:"<todo>"`
      Name    string   `xml:"<todo>"`
      UUID    string   `xml:"<todo>"`
      Version string   `xml:"<todo>"`
    }{}
    err := xml.Unmarshal(data, d)
    if err != nil {
      log.Println("ERROR: failed to read $info_type", err)
      return nil
    }
    info := &$info_type{}
    info.Name = d.Name
    info.UUID = d.UUID
    info.Version = d.Version
    return info
  }
  """
end

data = [
  ("ProcessInfo", "process"),
  ("FlowPropertyInfo", "flow property"),
  ("UnitGroupInfo", "unit group"),
  ("ContactInfo", "contact"),
  ("SourceInfo", "source")
]

for d in data
  println(gencode(d[1], d[2]))
end