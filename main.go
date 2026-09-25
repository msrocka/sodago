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
	tokens  *tokenAuth
	mutex   sync.Mutex
}

func main() {

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

	// create a default admin user when no user is configured yet
	if created, err := EnsureDefaultAdmin(args, config); err != nil {
		log.Fatalln("failed to create default admin user", err)
	} else if created {
		log.Println("INFO: default admin created (user: " +
			defaultAdminName + ", password: " + defaultAdminPassword + ")")
	}

	tokens, err := initTokenAuth(args)
	if err != nil {
		log.Fatalln("failed to init token auth", err)
	}

	server := server{
		config:  config,
		cookies: initCookieStore(args),
		tokens:  tokens,
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
	store := sessions.NewCookieStore(key)
	// sodago runs on plain HTTP, so the session cookie must not be marked as
	// secure; otherwise, standard clients would not send it back.
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return store
}
