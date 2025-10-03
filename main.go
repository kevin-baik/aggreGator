package main

import (
    "log"
    "github.com/kevin-baik/aggreGator/internal/config"
)

func main() {
    cfg, err := config.Read()
    if err != nil {
	log.Fatalf("Error reading config file: %v", err)
    }
   
    err = config.SetUser(cfg, "kevin")
    if err != nil {
	log.Fatalf("Error setting user config file: %v", err)
    }
    cfg, err = config.Read()
    if err != nil {
	log.Fatalf("Error reading config file: %v", err)
    }
}
