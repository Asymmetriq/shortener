package config

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

var (
	//go:embed migrations/*.sql
	migrationsDir     embed.FS
	migrationsDirName = "migrations"
)

func ConnectToDatabase(driver, dsn string) *sql.DB {
	if dsn == "" {
		return nil
	}
	con, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("establish db connection, %v", err)
	}

	migrateDatabase(con)
	return con
}

func migrateDatabase(db *sql.DB) {
	goose.SetBaseFS(migrationsDir)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrate db: %v", err)
	}

	if err := goose.Up(db, migrationsDirName); err != nil {
		log.Fatalf("migrate db up: %v", err)
	}
}
