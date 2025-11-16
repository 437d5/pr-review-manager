package db

import (
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migrator struct {
	conn *sqlx.DB
}

func NewMigrator(conn *sqlx.DB) *Migrator {
	return &Migrator{conn: conn}
}

func (db *Migrator) Migrate() error {
	slog.Debug("running migration")
	files, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db.conn.DB, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if err != migrate.ErrNoChange {
			slog.Debug("migration failed",
				slog.String("error", err.Error()),
			)
			return err
		}
		slog.Debug("migration did not change anything")
	}

	slog.Debug("migration finished")
	return nil
}
