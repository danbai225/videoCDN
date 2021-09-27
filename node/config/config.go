package config

import (
	"encoding/json"
	"io/ioutil"
)

var GlobalConfig Config
var Path = "./config.json"

type Config struct {
	ServerAddress string
	Token         string
	Port          int
	CacheDir      string
}

func LoadConfig() error {
	config := Config{}
	data, err := ioutil.ReadFile(Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.CacheDir == "" {
		config.CacheDir = "./cache"
	}
	GlobalConfig = config
	return nil
}
