package api

import "os"

type serverConfig struct {
	database databaseConfig
}

func loadServerConfig() *serverConfig {
	return &serverConfig{
		database: *loadDatabaseConfig(),
	}
}

type databaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func loadDatabaseConfig() *databaseConfig {
	return &databaseConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_DATABASE"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
}
