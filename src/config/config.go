package config

import (
	"encoding/json"
	"os"
)

const configPath = "~/.config/cfupdater/config.json"

type Config struct {
	// Mandatory
	Email  string `json:"email"`
	ApiKey string `json:"api_key"`
	Zone   string `json:"zone"`
	Name   string `json:"name"`

	// Optional
	Minute  int64 `json:"minute,omitempty"`
	Timeout int64 `json:"timeout,omitempty"`
}

var config = &Config{
	Minute:  15,
	Timeout: 30,
}

func InitConfig() error {
	homePath := os.Getenv("HOME")

	file, err := os.Open(homePath + configPath[1:])
	if nil != err {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(config)
}

func GetConfig() Config {
	return *config
}
