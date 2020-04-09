package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var db *DB
var cookieStore *sessions.CookieStore

func main() {

	_, err := newDataDir("data")
	if err != nil {
		log.Fatalln("failed to init data folder", err)
	}

	args := GetArgs()
	os.MkdirAll(args.DataDir, os.ModePerm)
	initCookieStore(args)

	// initialize the database
	log.Println("initialize data in:", args.DataDir)
	db, err = InitDB(args.DataDir)
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}
	db.RootDataStock()

	r := mux.NewRouter()
	registerRoutes(r, args)

	log.Println("Register shutdown routines")
	ossignals := make(chan os.Signal)
	signal.Notify(ossignals, syscall.SIGTERM)
	signal.Notify(ossignals, syscall.SIGINT)
	go func() {
		<-ossignals
		log.Println("Shutdown server")
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database", err)
		}
		os.Exit(0)
	}()

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
