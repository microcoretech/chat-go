// Copyright 2025 Mykhailo Bobrovskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
