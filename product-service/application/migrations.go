package application

import (
	"database/sql"
	"log"

	databaseinit "github.com/BernardN38/ebuy-server/product-service/application/databaseInit"
	"github.com/pressly/goose/v3"
)

func RunDatabaseMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	log.Println(embedMigrations, "helloooo")
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	if err := databaseinit.InitializeProductTypesTable(db); err != nil {
		return err
	}
	return nil
}
