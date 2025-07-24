package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDBConnection(migrationSource, dbHost, dbPort, dbUser, dbPassword, dbName, sslMode string) (*sqlx.DB, error) {
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)
	db, err := sqlx.Open("pgx", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}

	if err := runMigrations(migrationSource, dbDSN); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return db, nil
}

func runMigrations(migrationSource, dbDSN string) error {
	m, err := migrate.New(
		migrationSource,
		dbDSN)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == migrate.ErrNoChange {
		log.Println("Migrations: no change")
		return nil
	}
	log.Println("Migrations done")
	return nil
}
