package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Structs

// Config holds all information parsed from
// supplied config file.
type Config struct {
	Server   Server
	Settings Settings
	Session  Session
}

// Server holds all server information
// including hostname and port.
type Server struct {
	Addr string
	TLS  bool
}

// Settings holds all global parameters such
// as the number of threads and the seed to
// generate the involved IMAP commands.
type Settings struct {
	Threads  int
	Sessions int
	Seed     int64
	Throttle int
}

// Session holds all information about the
// length of one session.
type Session struct {
	MinLength int
	MaxLength int
}

// Functions

// LoadConfig decodes the config file and creates a
// Config object.
func LoadConfig(configFile string) (*Config, error) {

	conf := &Config{}

	// Parse values from TOML file into struct.
	if _, err := toml.DecodeFile(configFile, conf); err != nil {
		return nil, fmt.Errorf("failed to read in TOML config file at '%s' with: %v", configFile, err)
	}

	return conf, nil
}
