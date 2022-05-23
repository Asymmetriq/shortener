package db

import (
	"database/sql"
	"log"
)

func ConnectToPostgres(driver, dsn string) *sql.DB {
	con, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("establish db connection, %v", err)
	}
	return con
}
