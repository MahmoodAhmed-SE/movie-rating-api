package config

import (
	"log"
	constants "movie-rating-api-go/internals"
	"os"

	// "github.com/jackc/pgx"
	// "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	IP   string
	Port int
}

type Config struct {
	// PgxConfig    pgx.ConnConfig
	ServerConfig ServerConfig
}

func LoadConfig() Config {
	err := godotenv.Load()

	_, errUserEnv := os.LookupEnv(constants.EnvDBUser)
	_, errPassEnv := os.LookupEnv(constants.EnvDBPass)
	if !errUserEnv {
		log.Println("DB User .env variable error. Using System environment variables.")
		os.Getenv(constants.EnvDBUser)
	}
	if !errPassEnv {
		log.Println("DB Pass .env variable error. Using System environment variables.")
		os.Getenv(constants.EnvDBPass)
	}

	if err != nil {
		log.Println("Error loading .env variables: %v", err)
	}

	// postgresConfig := pgx.ConnConfig{
	// 	Config: pgconn.Config{

	// 	},
	// 	User:     os.Getenv(constants.EnvDBUser),
	// 	Password: os.Getenv(constants.EnvDBPass),
	// 	Host:     "localhost",
	// 	Database: "moviesdb",
	// 	Port:     5432,
	// }

	serverConfig := ServerConfig{
		IP:   "127.0.0.1",
		Port: 4040,
	}

	return Config{
		// PgxConfig:    postgresConfig,
		ServerConfig: serverConfig,
	}
}
