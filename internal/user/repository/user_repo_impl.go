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
	"fmt"
	"strings"

	"chat-go/internal/common/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/repository"
	"chat-go/internal/user/common"
)

type UserRepoImpl struct {
	db *sql.DB
}

var (
	userFieldsMapping = map[string]string{
		"id":        "id",
		"email":     "email",
		"username":  "username",
		"role":      "role",
		"firstName": "first_name",
		"lastName":  "last_name",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
)

func (r *UserRepoImpl) scan(rows *sql.Rows) ([]domain.User, error) {
	users := make([]domain.User, 0)

	for rows.Next() {
		var user domain.User

		fields := []any{
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Role,
			&user.FirstName,
			&user.LastName,
			&user.AboutMe,
			&user.Image.URL,
			&user.CreatedAt,
			&user.UpdatedAt,
		}

		err := rows.Scan(fields...)
		if err != nil {
			return nil, errors.NewDatabaseError(common.UserDomain, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepoImpl) buildFilter(filter domain.UserFilter) ([]any, []string) {
	values := make([]any, 0)
	where := make([]string, 0)

	if len(filter.IDs) > 0 {
		var params []string
		for _, id := range filter.IDs {
			values = append(values, id)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"u.id IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Emails) > 0 {
		var params []string
		for _, email := range filter.Emails {
			values = append(values, email)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"u.email IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.UserNames) > 0 {
		var params []string
		for _, userName := range filter.UserNames {
			values = append(values, userName)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"u.username IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Roles) > 0 {
		var params []string
		for _, role := range filter.Roles {
			values = append(values, role)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"role IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Search) > 0 {
		values = append(values, fmt.Sprintf("%%%s%%", filter.Search))
		where = append(where, fmt.Sprintf(
			"CONCAT(u.id, u.email, u.username, u.first_name, u.last_name) ILIKE $%d ", len(values)))
	}

	return values, where
}

func (r *UserRepoImpl) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	sql := fmt.Sprintf("SELECT %s FROM %s AS u WHERE id = $1", userFields, usersTableName)

	rows, err := r.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, errors.NewDatabaseError(common.UserDomain, err, "error on query user")
	}

	defer rows.Close()

	users, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return &users[0], nil
	}

	return nil, nil
}

func (r *UserRepoImpl) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	if filter == nil {
		filter = &domain.UserFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf("SELECT %s FROM %s AS u", userFields, usersTableName)

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	if filter.Sort != nil {
		query = fmt.Sprintf(`%s ORDER BY %s %s`,
			query,
			userFieldsMapping[filter.Sort.SortBy],
			filter.Sort.SortDir,
		)
	}

	if filter.Limit != nil {
		query = fmt.Sprintf(`%s LIMIT %d`, query, *filter.Limit)
	}

	if filter.Offset != nil {
		query = fmt.Sprintf(`%s OFFSET %d`, query, *filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, errors.NewDatabaseError(common.UserDomain, err, "error on query users")
	}

	defer rows.Close()

	users, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepoImpl) GetUsersCount(ctx context.Context, filter *domain.UserFilter) (uint64, error) {
	if filter == nil {
		filter = &domain.UserFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf("SELECT COUNT(*) AS count FROM %s AS u", usersTableName)

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, errors.NewDatabaseError(common.UserDomain, err, "error on query users count")
	}

	defer rows.Close()

	var count uint64

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, errors.NewDatabaseError(common.UserDomain, err, "error on scan users count")
		}
	}

	return count, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user domain.User, tx repository.Tx) (*domain.User, error) {
	values := []any{
		user.Email,
		user.Username,
		user.Role,
		user.FirstName,
		user.LastName,
	}

	query := fmt.Sprintf(`
		WITH u AS (
		    INSERT INTO %[1]s 
				(email, username, role, first_name, last_name)
			VALUES 
				($1, $2, $3, $4, $5)
			RETURNING *
		)
		SELECT %[2]s FROM u
	`, usersTableName, userFields)

	var (
		err  error
		rows *sql.Rows
	)

	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, values...)
	} else {
		rows, err = r.db.QueryContext(ctx, query, values...)
	}

	if err != nil {
		return nil, errors.NewDatabaseError(common.UserDomain, err)
	}

	defer rows.Close()

	users, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func NewUserRepoImpl(db *sql.DB) *UserRepoImpl {
	return &UserRepoImpl{
		db: db,
	}
}
