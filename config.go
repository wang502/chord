package chord

import (
	"crypto/sha1"
	"encoding/json"
	_ "errors"
	"hash"
	"io/ioutil"
	"log"
	_ "strings"
)

// Config represents configuration for a Chord node
type Config struct {
	Host          string `json:"Host"`
	HashFunc      hash.Hash
	HashBits      int `json:"NumBits"`
	NumNodes      int `json:"NumNodes"`
	NumSuccessors int `json:"NumSuccessors"`
}

// InitConfig initializes configuration from conf file
func InitConfig(confPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	log.Println(string(bytes))
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// DefaultConfig initializes a default configuration
func DefaultConfig(host string) *Config {
	return &Config{
		Host:          host,
		HashFunc:      sha1.New(),
		HashBits:      160,
		NumNodes:      8,
		NumSuccessors: 8,
	}
}
