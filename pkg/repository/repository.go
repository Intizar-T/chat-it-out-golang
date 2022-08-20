package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type BaseRepo interface {
	Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	CreateOne(ctx context.Context, query string, args interface{}) error
	Delete(ctx context.Context, query string, args ...interface{}) error
	Update(ctx context.Context, query string, args ...interface{}) error
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Repository struct {
	DB dbExecutor
}

// NewRepository creates a new repository
func NewRepository(executor dbExecutor) *Repository {
	return &Repository{
		DB: executor,
	}
}

func (repo *Repository) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return repo.DB.SelectContext(ctx, dest, query, args...)
}

func (repo *Repository) QueryOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return repo.DB.GetContext(ctx, dest, query, args...)

}

// CreateOne execute query to create one row.
// It returns error if rows inserted is less than 1.
func (repo *Repository) CreateOne(ctx context.Context, query string, args interface{}) error {

	var result sql.Result
	var err error
	result, err = repo.DB.NamedExecContext(ctx, query, args)

	if err != nil {
		return err
	}

	inserted, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if inserted < 1 {
		return fmt.Errorf("no record inserted after CreateOne with query: %s", query)
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, query string, args ...interface{}) error {
	result, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) Update(ctx context.Context, query string, args ...interface{}) error {
	result, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
