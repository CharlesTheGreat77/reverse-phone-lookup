package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
}

var (
	Token     string
	BotPrefix string

	config *Config
)

func ReadConfig(cwd string) error {
	fmt.Printf("[*] Reading config.json file..")
	file, err := os.ReadFile(fmt.Sprintf("%s/config.json", cwd))
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
