package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	Port     int
	DBFile   string
	Password string
}

var appConfig *AppConfig

func Load() *AppConfig {
	if appConfig != nil {
		return appConfig
	}

	port := 7540
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	dbFile := "scheduler.db"
	if customDB := os.Getenv("TODO_DB_FILE"); customDB != "" {
		dbFile = customDB
	}

	appConfig = &AppConfig{
		Port:     port,
		DBFile:   dbFile,
		Password: os.Getenv("TODO_PASSWORD"),
	}
	return appConfig
}

func Get() *AppConfig {
	if appConfig == nil {
		panic("configuration not loaded. Call config.Load() first")
	}
	return appConfig
}
