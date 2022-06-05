package database

import (
	"github.com/jmoiron/sqlx"
	"os"
)

func ConnectToDb() (*sqlx.DB, error) {
	connectionString := os.Getenv("CONNECTION_STRING")

	return sqlx.Connect("postgres", connectionString)
}
