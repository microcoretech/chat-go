package repository

import (
	"context"
	"database/sql"
)

type BaseRepo interface {
	Begin() (Tx, error)
	BeginContext(ctx context.Context) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type BaseRepoImpl struct {
	db *sql.DB
}

func (r *BaseRepoImpl) Begin() (Tx, error) {
	return r.db.Begin()
}

func (r *BaseRepoImpl) BeginContext(ctx context.Context) (Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func NewBaseRepoImpl(db *sql.DB) *BaseRepoImpl {
	return &BaseRepoImpl{db: db}
}
