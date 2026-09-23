package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type server struct {
	config  *Config
	dir     *datadir
	cookies *sessions.CookieStore
	mutex   sync.Mutex
}

func main() {

	// first check if a specific command was called
	osArgs := os.Args
	if len(osArgs) > 1 {
		cmd := osArgs[1]
		if cmd == "add-user" {
			AddUser()
			return
		}
	}

	args := ParseArgs()

	// ensure the data folder exists before anything tries to write into it
	if err := os.MkdirAll(args.DataDir(), os.ModePerm); err != nil {
		log.Fatalln("failed to create data folder", args.DataDir(), err)
	}

	config, err := ReadConfig(args)
	if err != nil {
		fmt.Println("ERROR: failed to read config", err)
		return
	}

	server := server{
		config:  config,
		cookies: initCookieStore(args),
	}
	dir, err := newDataDir(args.DataDir())
	if err != nil {
		log.Fatalln("failed to init data folder", err)
	}
	server.dir = dir

	r := mux.NewRouter()
	server.registerRoutes(r)

	log.Println("Starting server at port:", args.Port())
	http.ListenAndServe(":"+args.Port(), r)
}

func initCookieStore(args Args) *sessions.CookieStore {
	log.Println("Init cookie store ...")
	keyPath := filepath.Join(args.DataDir(), "cookie_auth.key")
	_, err := os.Stat(keyPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("Cannot access cookie key at", keyPath, err)
	}
	var key []byte
	if os.IsNotExist(err) {
		key = securecookie.GenerateRandomKey(32)
		err = os.WriteFile(keyPath, key, os.ModePerm)
		if err != nil {
			log.Fatalln("Failed to save", keyPath, ": ", err)
		}
	} else {
		key, err = os.ReadFile(keyPath)
		if err != nil {
			log.Fatalln("Failed to read", keyPath, ": ", err)
		}
	}
	return sessions.NewCookieStore(key)
}
