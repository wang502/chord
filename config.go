package chord

import (
    "encoding/json"
    "crypto/sha1"
    "hash"
    _ "errors"
    _ "strings"
    "io/ioutil"
    "log"
)

// configuration for a Chord node
type Config struct {
    Host string `json:"Host"`
    HashFunc hash.Hash
    HashBits int `json:"NumBits"`
    NumNodes int `json:"NumNodes"`
    NumSuccessors int `json:"NumSuccessors"`
}

// initialize configuration from conf file
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

// initialize default configuration
func DefaultConfig(host string) (*Config) {
    return &Config{
        Host          :   host,
        HashFunc      :   sha1.New(),
        HashBits      :   160,
        NumNodes      :     8,
        NumSuccessors :     8,
    }
}
