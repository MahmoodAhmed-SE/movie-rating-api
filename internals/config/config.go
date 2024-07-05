package config

import (
	"log"
	constants "movie-rating-api-go/internals"
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

	_, errUserEnv:= os.LookupEnv(constants.EnvDBUser)
	_, errPassEnv:= os.LookupEnv(constants.EnvDBPass)
	if !errUserEnv {
		log.Fatalf("DB User environment variable error: %v", errUserEnv)
	}
	if !errPassEnv {
		log.Fatalf("DB Pass environment variable error: %v", errPassEnv)
	}

	
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	postgresConfig := pgx.ConnConfig{
		User:     os.Getenv(constants.EnvDBUser),
		Password: os.Getenv(constants.EnvDBPass),
		Host:     "localhost",
		Database: "moviesdb",
		Port:     5432,
	}

	serverConfig := ServerConfig{
		IP:   "127.0.0.1",
		Port: 4040,
	}

	return Config{
		PgxConfig:    postgresConfig,
		ServerConfig: serverConfig,
	}
}
