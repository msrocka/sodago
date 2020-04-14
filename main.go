package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var cookieStore *sessions.CookieStore

type server struct {
	dir   *datadir
	mutex sync.Mutex
}

func main() {

	server := server{}
	dir, err := newDataDir("data")
	if err != nil {
		log.Fatalln("failed to init data folder", err)
	}
	server.dir = dir

	args := GetArgs()
	os.MkdirAll(args.DataDir, os.ModePerm)
	initCookieStore(args)

	r := mux.NewRouter()
	server.registerRoutes(r, args)

	log.Println("Starting server at port:", args.Port)
	http.ListenAndServe(":"+args.Port, r)
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
