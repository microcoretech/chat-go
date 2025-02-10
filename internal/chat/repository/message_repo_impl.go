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

	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/chat/domain"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/common/repository"
)

type MessageRepoImpl struct {
	db *sql.DB
}

var (
	messageFieldsMapping = map[string]string{
		"id":        "id",
		"status":    "status",
		"chatId":    "chat_id",
		"createdBy": "created_by",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
)

func (r *MessageRepoImpl) scan(rows *sql.Rows) ([]domain.Message, error) {
	if rows == nil {
		return nil, nil
	}

	messages := make([]domain.Message, 0)

	for rows.Next() {
		var message domain.Message

		var fields = []any{
			&message.ID,
			&message.Text,
			&message.Status,
			&message.ChatID,
			&message.CreatedBy,
			&message.CreatedAt,
			&message.UpdatedAt,
		}

		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepoImpl) buildFilter(filter domain.MessageFilter) ([]any, []string) {
	values := make([]any, 0)
	where := make([]string, 0)

	if len(filter.IDs) > 0 {
		var params []string
		for _, id := range filter.IDs {
			values = append(values, id)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"m.id IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.ChatIDs) > 0 {
		var params []string
		for _, chatID := range filter.ChatIDs {
			values = append(values, chatID)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"m.chat_id IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.CreatedByIDs) > 0 {
		var params []string
		for _, createdByID := range filter.CreatedByIDs {
			values = append(values, createdByID)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"m.chat_id IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Statuses) > 0 {
		var params []string
		for _, status := range filter.Statuses {
			values = append(values, status)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"m.chat_id IN (%s) ", strings.Join(params, ",")))
	}

	return values, where
}

func (r *MessageRepoImpl) GetMessages(ctx context.Context, filter *domain.MessageFilter) ([]domain.Message, error) {
	if filter == nil {
		filter = &domain.MessageFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s AS m
	`, messageFields, messageTableName)

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	if filter.Sort != nil {
		query = fmt.Sprintf(`%s ORDER BY %s %s`,
			query,
			messageFieldsMapping[filter.Sort.SortBy],
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
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	messages, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepoImpl) GetMessagesCount(ctx context.Context, filter *domain.MessageFilter) (uint64, error) {
	if filter == nil {
		filter = &domain.MessageFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf("SELECT COUNT(*) AS count FROM %s AS m", messageTableName)

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	var count uint64

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, errors.NewDatabaseError(common.ChatDomain, err)
		}
	}

	return count, nil
}

func (r *MessageRepoImpl) CreateMessage(ctx context.Context, message domain.Message, tx repository.Tx) (*domain.Message, error) {
	values := []any{
		message.Text,
		message.ChatID,
		message.CreatedBy,
	}

	query := fmt.Sprintf(`
		WITH %[1]s AS (
		    INSERT INTO %[1]s (
				text,
				chat_id,
				created_by
			)
			VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT %[2]s
		FROM %[1]s AS m
	`,
		messageTableName,
		messageFields,
	)

	var (
		rows *sql.Rows
		err  error
	)

	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, values...)
	} else {
		rows, err = r.db.QueryContext(ctx, query, values...)
	}

	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	messages, err := r.scan(rows)
	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[0], nil
}

func (r *MessageRepoImpl) UpdateMessageStatus(
	ctx context.Context,
	messageIDs []uint64,
	messageStatus domain.MessageStatus,
	tx repository.Tx,
) error {
	if len(messageIDs) == 0 {
		return nil
	}

	var values []any

	query := fmt.Sprintf("UPDATE %s", messageTableName)

	values = append(values, messageStatus)
	query = fmt.Sprintf("%s SET status = %s", query, fmt.Sprintf("$%d", len(values)))

	var placeholders []string
	for _, id := range messageIDs {
		values = append(values, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	query = fmt.Sprintf("%s WHERE id IN (%s)", query, strings.Join(placeholders, ", "))

	values = append(values, messageStatus)
	query = fmt.Sprintf("%s AND status < %s", query, fmt.Sprintf("$%d", len(values)))

	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, values...)
	} else {
		_, err = r.db.QueryContext(ctx, query, values...)
	}

	if err != nil {
		return err
	}

	return nil
}

func NewMessageRepoImpl(db *sql.DB) *MessageRepoImpl {
	return &MessageRepoImpl{db: db}
}
