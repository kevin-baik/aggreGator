package config

import (
    "fmt"
    "os"
    "encoding/json"
    "path/filepath"
)

const configFileName = ".gatorconfig.json"

// .gatorconfig JSON configuration
type Config struct {
    DBUrl		string `json:"db_url"`
    CurrentUserName	string `json:"current_user_name"`
}

func (c Config) String() string {
    return fmt.Sprintf("Config {\n\tdbUrl: %v\n\tusername: %v\n", c.DBUrl, c.CurrentUserName)
}

func (cfg *Config) SetUser(username string) error {
    cfg.CurrentUserName = username
    return write(*cfg)
}

func Read() (Config, error) {
    configPath, err := getConfigFilePath()
    if err != nil {
	return Config{}, err
    }
    file, err := os.Open(configPath)
    if err != nil {
	return Config{}, err
    }
    defer file.Close()
    
    var cfg Config
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&cfg)
    if err != nil {
	return Config{}, err
    }
    fmt.Println(cfg)
    return cfg, nil
}

func getConfigFilePath() (string, error) {
    homepath, err := os.UserHomeDir()
    if err != nil {
	return "", err
    }
    configPath := filepath.Join(homepath, configFileName)
    return configPath, nil
}

func write(cfg Config) error {
    configPath, err := getConfigFilePath()
    if err != nil {
	return err
    }

    file, err := os.Create(configPath)
    if err != nil {
	return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(cfg)
    if err != nil {
	return err
    }
    
    fmt.Printf("Configuration successfully written to: %s\n", configPath)
    fmt.Println(cfg)
    return nil
}
