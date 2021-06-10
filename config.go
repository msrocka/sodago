package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
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

// ReadConfig reads the configuration file from the data folder.
func ReadConfig(args Args) (*Config, error) {
	config := &Config{}
	path := filepath.Join(args.DataDir(), "config.json")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, config); err != nil {
		return nil, err
	} else {
		return config, nil
	}
}

func WriteConfig(args Args, config *Config) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(args.DataDir(), "config.json")
	return ioutil.WriteFile(path, bytes, os.ModePerm)
}

func AddUser() {
	// parse and check the arguments
	args := ParseArgs()
	name := args["-name"]
	pw := args["-password"]
	if name == "" || pw == "" {
		fmt.Println("ERROR: no user or password given")
		fmt.Println("To add a user the command should be:")
		fmt.Println("  sodago add-user -name [USER_NAME] -password [PASSWORD]")
		return
	}

	config, err := ReadConfig(args)
	if err != nil {
		fmt.Println("ERROR: failed to read configuration file:", err)
		return
	}

	// check that the user does not exist yet
	for _, user := range config.Users {
		if user.Name == name {
			fmt.Println("ERROR: a user", name, "already exists")
			return
		}
	}

	// update the configuration
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR: failed to hash password", err)
		return
	}
	config.Users = append(config.Users, User{
		Name: name,
		Hash: base64.StdEncoding.EncodeToString(hash),
	})
	if err = WriteConfig(args, config); err != nil {
		fmt.Println("ERROR: failed to write configuration file:", err)
	}
}
