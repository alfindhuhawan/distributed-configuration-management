package config

import (
	"encoding/json"
	"log"
	"os"
)

func ReadConfig() *Config {

	bytes, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg Config

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err.Error())
	}

	return &cfg
}

func ReadAgentCache(filename string) *AgentCache {

	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg AgentCache

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err.Error())
	}

	return &cfg
}
