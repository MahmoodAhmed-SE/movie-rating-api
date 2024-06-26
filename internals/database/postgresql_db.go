package database

import (
	"github.com/jackc/pgx"
)

var conn *pgx.Conn

func ConnectPostgres(config pgx.ConnConfig) error {
	var connErr error
	conn, connErr = pgx.Connect(config)

	if connErr != nil {
		return connErr
	}

	return nil
}

func GetConn() *pgx.Conn {
	return conn
}
