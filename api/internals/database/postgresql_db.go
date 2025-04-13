package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func ConnectPostgres(config pgx.ConnConfig) error {
	var connErr error
	connStr := os.Getenv("DATABASE_URL")
	conn, connErr = pgx.Connect(context.Background(), connStr)

	if connErr != nil {
		return connErr
	}

	return nil
}

func GetConn() *pgx.Conn {
	return conn
}
