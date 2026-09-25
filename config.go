package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// User contains the data of a registered user with user name and password.
// Note: this application is for testing only, therefore the password is stored
// as plain text in the configuration file.
type User struct {
	Name     string   `json:"user"`
	Password string   `json:"password"`
	Roles    []string `json:"roles,omitempty"`
	// Tokens are pre-configured authentication tokens of this user. They are
	// an alternative to the tokens that are generated via the service API.
	Tokens []string `json:"tokens,omitempty"`
}

// Config holds the configuration of the users and data stocks.
type Config struct {
	Users []User `json:"users"`
}

// GetUser returns the user with the given name from the configuration. Returns
// nil if there is no such user.
func (config *Config) GetUser(name string) *User {
	if config == nil || name == "" {
		return nil
	}
	lowerName := strings.ToLower(strings.TrimSpace(name))
	for i := range config.Users {
		user := &config.Users[i]
		if strings.ToLower(user.Name) == lowerName {
			return user
		}
	}
	return nil
}

// GetUserByToken returns the user that has the given token in its configured
// token list. Returns nil if there is no such user.
func (config *Config) GetUserByToken(token string) *User {
	if config == nil || token == "" {
		return nil
	}
	for i := range config.Users {
		for _, configured := range config.Users[i].Tokens {
			if configured == token {
				return &config.Users[i]
			}
		}
	}
	return nil
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

	bytes, err := os.ReadFile(path)
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
	return os.WriteFile(path, bytes, os.ModePerm)
}

// saveConfig persists the current configuration in the data folder.
func (s *server) saveConfig() error {
	return WriteConfig(Args{"-data": s.dir.root}, s.config)
}

// defaultAdminName and defaultAdminPassword define the credentials of the
// admin user that is created when the configuration does not contain any user
// yet.
const (
	defaultAdminName     = "admin"
	defaultAdminPassword = "default"
)

// EnsureDefaultAdmin creates a default admin user and persists the
// configuration when no user is defined yet. It returns true when a default
// admin was created.
func EnsureDefaultAdmin(args Args, config *Config) (bool, error) {
	if len(config.Users) > 0 {
		return false, nil
	}

	config.Users = append(config.Users, User{
		Name:     defaultAdminName,
		Password: defaultAdminPassword,
		Roles:    []string{"admin"},
	})
	if err := WriteConfig(args, config); err != nil {
		return false, err
	}
	return true, nil
}
