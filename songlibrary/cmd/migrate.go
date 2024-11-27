package cmd

import (
	stdsql "database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"path/filepath"
)

func Migrate(config Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	fmt.Println(connStr)

	db, err := stdsql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := filepath.Abs("./songlibrary/db/migrations")
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", migrations),
		config.DBName, driver)
	if err != nil {
		return err
	}

	err = migration.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
