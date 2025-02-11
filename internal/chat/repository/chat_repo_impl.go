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

	"chat-go/internal/chat/common"
	"chat-go/internal/chat/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/repository"
)

type ChatRepoImpl struct {
	db *sql.DB
}

var (
	chatFieldsMapping = map[string]string{
		"id":        "id",
		"name":      "name",
		"type":      "type",
		"createdBy": "created_by",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
)

const (
	chatFields        = `c.id, c.name, c.type, c.image_url, c.created_by, c.created_at, c.updated_at`
	lastMessageFields = `
		(
			SELECT
				JSONB_BUILD_OBJECT(
					'id', m.id,
					'text', m.text,
					'status', m.status,
					'chatId', m.chat_id,
					'createdBy', m.created_by,
					'createdAt', CAST(m.created_at as timestamp) AT time zone 'UTC',
					'updatedAt', CAST(m.updated_at AS timestamp) AT time zone 'UTC'
				)
			FROM messages AS m WHERE m.chat_id = c.id ORDER BY m.updated_at DESC LIMIT 1
		) as last_message
	`
	userChatFields = `COALESCE(
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'userId', uc.user_id,
				'chatId', uc.chat_id
			)
		) FILTER (WHERE uc.user_id IS NOT NULL), '[]'::JSON) AS user_chats
	`
)

func (r *ChatRepoImpl) scan(rows *sql.Rows) ([]domain.Chat, error) {
	if rows == nil {
		return nil, nil
	}

	chats := make([]domain.Chat, 0)

	for rows.Next() {
		var chat domain.Chat
		var lastMessage *messageDto

		var fields = []any{
			&chat.ID,
			&chat.Name,
			&chat.Type,
			&chat.Image.URL,
			&chat.CreatedBy,
			&chat.CreatedAt,
			&chat.UpdatedAt,
			&lastMessage,
			(*userChatsDto)(&chat.UserChats),
		}

		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}

		chat.LastMessage = (*domain.Message)(lastMessage)
		chats = append(chats, chat)
	}

	return chats, nil
}

func (r *ChatRepoImpl) buildChatFields() string {
	return fmt.Sprintf(`%s, %s, %s`, chatFields, lastMessageFields, userChatFields)
}

func (r *ChatRepoImpl) buildFrom() string {
	from := fmt.Sprintf(
		`%s AS c LEFT JOIN %s AS uc ON c.id = uc.chat_id`,
		chatTableName,
		userChatTableName,
	)

	return from
}

func (r *ChatRepoImpl) buildFilter(filter domain.ChatFilter) ([]any, []string) {
	values := make([]any, 0)
	where := make([]string, 0)

	if len(filter.IDs) > 0 {
		var params []string
		for _, id := range filter.IDs {
			values = append(values, id)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"c.id IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.CreatedByIDs) > 0 {
		var params []string
		for _, createdByID := range filter.CreatedByIDs {
			values = append(values, createdByID)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"c.created_by IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Types) > 0 {
		var params []string
		for _, t := range filter.Types {
			values = append(values, t)
			params = append(params, fmt.Sprintf("$%d", len(values)))
		}
		where = append(where, fmt.Sprintf(
			"c.type IN (%s) ", strings.Join(params, ",")))
	}

	if len(filter.Search) > 0 {
		values = append(values, fmt.Sprintf("%%%s%%", filter.Search))
		where = append(where, fmt.Sprintf(
			"CONCAT(c.id, c.name) ILIKE $%d ", len(values)))
	}

	return values, where
}

func (r *ChatRepoImpl) buildGroupBy() string {
	return "c.id"
}

func (r *ChatRepoImpl) GetChat(ctx context.Context, id uint64) (*domain.Chat, error) {
	query := fmt.Sprintf(
		`SELECT %s 
			FROM %s
			WHERE c.id = $1`,
		r.buildChatFields(),
		r.buildFrom(),
	)

	query = fmt.Sprintf(`%s 
		GROUP BY %s`, query, r.buildGroupBy())

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	chats, err := r.scan(rows)
	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	if len(chats) == 0 {
		return nil, nil
	}

	return &chats[0], nil
}

func (r *ChatRepoImpl) GetChats(ctx context.Context, filter *domain.ChatFilter) ([]domain.Chat, error) {
	if filter == nil {
		filter = &domain.ChatFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf("SELECT %s FROM %s", r.buildChatFields(), r.buildFrom())

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	query = fmt.Sprintf(`%s 
		GROUP BY %s`, query, r.buildGroupBy())

	query = fmt.Sprintf(`
		WITH r AS (%s)
		SELECT * FROM r`, query)

	if filter.Sort != nil {
		query = fmt.Sprintf(`%s
		ORDER BY `, query)
		query = fmt.Sprintf(`%s %s %s`,
			query,
			chatFieldsMapping[filter.Sort.SortBy],
			filter.Sort.SortDir,
		)
	} else {
		query = fmt.Sprintf(`%s 
		ORDER BY (COALESCE((last_message->'createdAt')::VARCHAR, '""'::VARCHAR)) DESC, updated_at DESC`, query)
	}

	if filter.Limit != nil {
		query = fmt.Sprintf(`%s 
		LIMIT %d`, query, *filter.Limit)
	}

	if filter.Offset != nil {
		query = fmt.Sprintf(`%s 
		OFFSET %d`, query, *filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	chats, err := r.scan(rows)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepoImpl) GetChatsCount(ctx context.Context, filter *domain.ChatFilter) (uint64, error) {
	if filter == nil {
		filter = &domain.ChatFilter{}
	}

	values, where := r.buildFilter(*filter)

	query := fmt.Sprintf("SELECT COUNT(*) AS count FROM %s AS c", chatTableName)

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, errors.NewDatabaseError(common.ChatDomain, err, "error on query chats count")
	}

	defer rows.Close()

	var count uint64

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, errors.NewDatabaseError(common.ChatDomain, err, "error on scan chats count")
		}
	}

	return count, nil
}

func (r *ChatRepoImpl) CreateChat(ctx context.Context, chat domain.Chat, tx repository.Tx) (*domain.Chat, error) {
	var imageURL string

	var err error

	values := []any{
		chat.Name,
		chat.Type,
		imageURL,
		chat.CreatedBy,
	}

	query := fmt.Sprintf(`
		WITH %[1]s AS (
		    INSERT INTO %[1]s (
				name,
				type,
		        image_url,
				created_by
			)
			VALUES ($1, $2, $3, $4)
			RETURNING *
		)
		SELECT %[2]s
		FROM %[3]s
	`,
		chatTableName,
		r.buildChatFields(),
		r.buildFrom(),
	)

	var rows *sql.Rows

	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, values...)
	} else {
		rows, err = r.db.QueryContext(ctx, query, values...)
	}

	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	chats, err := r.scan(rows)
	if err != nil {
		return nil, errors.NewDatabaseError(common.ChatDomain, err)
	}

	if len(chats) == 0 {
		return nil, nil
	}

	return &chats[0], nil
}

func (r *ChatRepoImpl) UpdateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ChatRepoImpl) DeleteChat(ctx context.Context, id uint64) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func NewChatRepoImpl(db *sql.DB) *ChatRepoImpl {
	return &ChatRepoImpl{db: db}
}
