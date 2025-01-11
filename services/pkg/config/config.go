package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	MongoDB struct {
		ConnectionString string `json:"ConnectionString"`
		Database         string `json:"Database"`
		Collections      struct {
			Users string `json:"Users"`
		} `json:"Collections"`
	} `json:"MongoDB"`
	PostgreSQL struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"postgres"`
	Server struct {
		Port string `json:"Port"`
	} `json:"Server"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Try current directory first
	configPath := "config/appSettings.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If not found, try parent directory
		configPath = filepath.Join("..", "config", "appSettings.json")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// Environment variables override config file
	if mongoURI := os.Getenv("MONGODB_URI"); mongoURI != "" {
		config.MongoDB.ConnectionString = mongoURI
	}

	return &config, nil
}
