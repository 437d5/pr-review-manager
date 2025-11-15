package db

import (
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB
}

func New(address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		slog.Error("connection problem",
			slog.String("address", address),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &DB{conn: db}, nil
}
