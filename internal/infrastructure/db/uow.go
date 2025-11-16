package db

import (
	"context"
	"fmt"

	"github.com/437d5/pr-review-manager/internal/domain/repositories"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Teams() repositories.TeamRepository {
	if u.tx != nil {
		return NewTeamRepository(u.tx)
	}
	return NewTeamRepository(u.db)
}

func (u *UnitOfWork) Users() repositories.UserRepository {
	if u.tx != nil {
		return NewUserRepository(u.tx)
	}
	return NewUserRepository(u.db)
}

func (u *UnitOfWork) PR() repositories.PullRequestRepository {
	if u.tx != nil {
		return NewPullRequestRepository(u.tx)
	}
	return NewPullRequestRepository(u.db)
}

func (u *UnitOfWork) Begin(ctx context.Context) error {
	if u.tx != nil {
		return fmt.Errorf("transaction already started")
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}

	u.tx = tx
	return nil
}

func (u *UnitOfWork) Commit() error {
	if u.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}

	// will process err in service logic
	err := u.tx.Commit()
	u.tx = nil
	return err
}

func (u *UnitOfWork) Rollback() error {
	if u.tx == nil {
		return nil
	}

	err := u.tx.Rollback()
	u.tx = nil
	return err
}

func (u *UnitOfWork) Close() error {
	return u.Rollback()
}
