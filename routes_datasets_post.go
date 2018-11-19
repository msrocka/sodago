package main

import (
	"net/http"
)

func postFlow(data []byte, stock *DataStock, w http.ResponseWriter) {
	info := ReadFlowInfo(data)
	if info == nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}
	content := db.Content(stock)
	if content.ContainsFlow(info) {
		http.Error(w, "Fow "+info.UUID+" "+info.Version+
			"already exists", http.StatusBadRequest)
		return
	}
	v := ParseVersion(info.Version).String() // standard format
	key := stock.ID + "/Flow/" + info.UUID + "/" + v
	db.Put(DataSetBucket, key, data)
	content.Flows = append(content.Flows, *info)
	db.UpdateContent(stock, content)
	ServeXML(stock, w)
}
