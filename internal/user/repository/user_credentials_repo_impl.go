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

	"chat-go/internal/common/errors"
	common_repo "chat-go/internal/common/repository"
	"chat-go/internal/user/common"
	"chat-go/internal/user/domain"
)

type UserCredentialsRepoImpl struct {
	db *sql.DB
}

func (r *UserCredentialsRepoImpl) scan(rows *sql.Rows) ([]domain.UserCredentials, error) {
	userCredentialsList := make([]domain.UserCredentials, 0)

	for rows.Next() {
		var userCredentials domain.UserCredentials

		fields := []any{
			&userCredentials.UserID,
			&userCredentials.Password,
			&userCredentials.User.ID,
			&userCredentials.User.Email,
			&userCredentials.User.Username,
			&userCredentials.User.Role,
			&userCredentials.User.FirstName,
			&userCredentials.User.LastName,
			&userCredentials.User.AboutMe,
			&userCredentials.User.Image.URL,
			&userCredentials.User.CreatedAt,
			&userCredentials.User.UpdatedAt,
		}

		err := rows.Scan(fields...)
		if err != nil {
			return nil, errors.NewDatabaseError(common.UserDomain, err)
		}

		userCredentialsList = append(userCredentialsList, userCredentials)
	}

	return userCredentialsList, nil
}

func (r *UserCredentialsRepoImpl) GetUserCredentialsByUsername(ctx context.Context, username string) (*domain.UserCredentials, error) {
	sql := fmt.Sprintf(
		`
			SELECT %s, %s FROM %s AS uc
			INNER JOIN %s AS u ON uc.user_id = u.id
			WHERE u.email = $1 OR u.username = $1
		`,
		userCredentialsFields,
		userFields,
		userCredentialsTableName,
		usersTableName,
	)

	rows, err := r.db.QueryContext(ctx, sql, username)
	if err != nil {
		return nil, errors.NewDatabaseError(common.UserDomain, err)
	}

	defer rows.Close()

	userCredentials, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	if len(userCredentials) > 0 {
		return &userCredentials[0], nil
	}

	return nil, nil
}

func (r *UserCredentialsRepoImpl) CreateUserCredentials(
	ctx context.Context,
	userCredentials domain.UserCredentials,
	tx common_repo.Tx,
) (*domain.UserCredentials, error) {
	values := []any{userCredentials.UserID, userCredentials.Password}

	query := fmt.Sprintf(`
		WITH uc AS (
			INSERT INTO %[1]s 
				(user_id, password)
			VALUES
				($1, $2)
			RETURNING *
		)	
		SELECT %[2]s, %[3]s FROM uc
		INNER JOIN %[4]s AS u ON uc.user_id = u.id
	`,
		userCredentialsTableName,
		userCredentialsFields,
		userFields,
		usersTableName)

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

	userCredentialsSlice, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	if len(userCredentialsSlice) == 0 {
		return nil, nil
	}

	return &userCredentialsSlice[0], nil
}

func NewUserCredentialsRepoImpl(db *sql.DB) *UserCredentialsRepoImpl {
	return &UserCredentialsRepoImpl{db: db}
}
