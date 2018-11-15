package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Context contains the application data
type Context struct {
	DataStocks []*DataStock
}

var db *DB
var cookieStore *sessions.CookieStore

func main() {
	args := GetArgs()
	os.MkdirAll(args.DataDir, os.ModePerm)
	initCookieStore(args)

	// initialize the database
	var err error
	log.Println("initialize data in:", args.DataDir)
	db, err = InitDB(args.DataDir)
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}

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

func initCookieStore(args *Args) {
	log.Println("Init cookie store ...")
	keyPath := filepath.Join(args.DataDir, "cookie_auth.key")
	_, err := os.Stat(keyPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("Cannot access cookie key at", keyPath, err)
	}
	var key []byte
	if os.IsNotExist(err) {
		key = securecookie.GenerateRandomKey(32)
		err = ioutil.WriteFile(keyPath, key, os.ModePerm)
		if err != nil {
			log.Fatalln("Failed to save", keyPath, ": ", err)
		}
	} else {
		key, err = ioutil.ReadFile(keyPath)
		if err != nil {
			log.Fatalln("Failed to read", keyPath, ": ", err)
		}
	}
	cookieStore = sessions.NewCookieStore(key)
}
