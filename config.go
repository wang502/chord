package chord

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
)

// Config represents configuration for a Chord node
type Config struct {
	Host     string `json:"Host"`
	HashFunc hash.Hash
	HashBits int `json:"NumBits"`
	NumNodes int `json:"NumNodes"`
}

// InitConfig initializes configuration from conf file
func InitConfig(confPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	log.Println(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("init config failed: %s", err)
	}

	config := Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("init config failed: %s", err)
	}
	return &config, nil
}

// DefaultConfig initializes a default configuration
func DefaultConfig(host string) *Config {
	return &Config{
		Host:     host,
		HashFunc: sha1.New(),
		HashBits: 3,
		NumNodes: 8,
	}
}
