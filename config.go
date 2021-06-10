package main

import (
	"log"
	"os"
)

// User contains the data of a registered user with user name and hashed
// password.
type User struct {
	Name string `json:"user"`
	Hash string `json:"hash"`
}

// Config holds the configuration of the users and data stocks.
type Config struct {
	Users []User `json:"users"`
}

//
func AddUser() {
	args := os.Args
	if len(args) < 4 {
		log.Fatalln("")
	}
}
