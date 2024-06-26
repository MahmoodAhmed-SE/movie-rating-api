package config

import (
	"log"
	"os"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	IP   string
	Port int
}

type Config struct {
	PgxConfig    pgx.ConnConfig
	ServerConfig ServerConfig
}

func LoadConfig() Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	postgresConfig := pgx.ConnConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     "localhost",
		Database: "moviesdb",
		Port:     5432,
	}

	serverConfig := ServerConfig{
		IP:   "127.0.0.1",
		Port: 8080,
	}

	return Config{
		PgxConfig:    postgresConfig,
		ServerConfig: serverConfig,
	}
}
