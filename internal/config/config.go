package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}

func write(cfg Config) error {
	//Get configFilePath
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Error Getting Config File Path: %v\n", err)
		return err
	}

	//Marshal JSON Data
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	//Write JOSN Data to configFilePath
	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error Writing JSON Data to %s: %v", configFileName, err)
		return err
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Error Getting Config File Path: %v\n", err)
		return Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Error Reading %s File: %v\n", configFileName, err)
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error Unmarshaling %s: %v\n", configFileName, err)
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(userName string) error {
	//Set UserName
	c.CurrentUserName = userName

	//Write config struct to config file
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
