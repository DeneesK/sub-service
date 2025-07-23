package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db *sql.DB
}

func InitDBConnection(migrationSource, dbHost, dbPort, dbUser, dbPassword, dbName, sslMode string) *PostgresDB {
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	runMigrations(migrationSource, dbDSN)
	p := &PostgresDB{db: db}
	return p
}

func runMigrations(migrationSource, dbDSN string) {
	m, err := migrate.New(
		migrationSource,
		dbDSN)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("Migrations: no change")
	} else {
		log.Println("Migrations done")
	}
}
