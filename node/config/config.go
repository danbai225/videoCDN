package config

import (
	"encoding/json"
	"io/ioutil"
)

var GlobalConfig Config

type Config struct {
	ServerAddress string
	Token         string
}

func LoadConfig(path ...string) error {
	config := Config{}
	configPath := "./config.json"
	if len(path) > 0 && path[0] != "" {
		configPath = path[0]
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	GlobalConfig = config
	return nil
}
