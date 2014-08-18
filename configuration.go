package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	GithubToken string
}

func (configuration *Configuration) Persist() error {
	bdata, marshallErr := json.Marshal(configuration)
	if marshallErr != nil {
		return marshallErr
	}
	os.Remove("config.json")
	file, createErr := os.Create("config.json")
	if createErr != nil {
		return createErr
	}
	file.WriteString(string(bdata))
	return nil
}

func NewConfiguration(token string) *Configuration {
	return &Configuration{token}
}

func OpenConfiguration() (*Configuration, error) {
	if _, err := os.Stat("config.json"); err != nil {
		return nil, err
	}

	file, openErr := os.Open("config.json")
	if openErr != nil {
		return nil, openErr
	}

	var configuration = &Configuration{}
	decoder := json.NewDecoder(file)
	decodeErr := decoder.Decode(&configuration)
	return configuration, decodeErr
}
