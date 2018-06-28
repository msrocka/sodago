package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Context contains the application data
type Context struct {
	DataStocks []*DataStock
}

func main() {

	context := &Context{
		DataStocks: InitStocks(),
	}

	router := mux.NewRouter()

	// data stocs
	router.Methods("GET").Path("/resource/datastocks").
		HandlerFunc(GetDataStocksHandler(context))

	// profiles
	router.Methods("GET").Path("/resource/profiles").
		HandlerFunc(GetProfileDescriptors)
	router.Methods("GET").Path("/resource/profiles/").
		HandlerFunc(GetProfileDescriptors)
	router.Methods("GET").Path("/resource/profiles/{id}").
		HandlerFunc(GetProfile)
	http.ListenAndServe(":8080", router)
}
