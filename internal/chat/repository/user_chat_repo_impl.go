package repository

import (
	"chat/internal/chat/common"
	"chat/internal/chat/domain"
	"chat/internal/common/errors"
	"chat/internal/common/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type UserChatRepoImpl struct {
	db *sql.DB
}

func (r *UserChatRepoImpl) CreateUserChats(ctx context.Context, userChats []domain.UserChat, tx repository.Tx) error {
	if len(userChats) == 0 {
		return nil
	}

	var (
		placeholders []string
		values       []interface{}
	)

	const colsNum = 2

	for i, userChat := range userChats {
		var indexes []any

		for j := 1; j <= colsNum; j++ {
			indexes = append(indexes, i*colsNum+j)
		}

		placeholder := fmt.Sprintf("($%d,$%d)", indexes...)

		placeholders = append(placeholders, placeholder)

		values = append(values,
			userChat.UserID,
			userChat.ChatID,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
		VALUES %s
		ON CONFLICT DO NOTHING
	`,
		userChatTableName,
		strings.Join(placeholders, ","),
	)

	var (
		rows *sql.Rows
		err  error
	)

	if tx != nil {
		rows, err = tx.Query(query, values...)
	} else {
		rows, err = r.db.Query(query, values...)
	}

	if err != nil {
		return errors.NewDatabaseError(common.ChatDomain, err)
	}

	defer rows.Close()

	return nil
}

func (r *UserChatRepoImpl) DeleteUserChats(ctx context.Context, userChat []domain.UserChat, tx repository.Tx) error {
	//TODO implement me
	panic("implement me")
}

func NewUserChatRepoImpl(db *sql.DB) *UserChatRepoImpl {
	return &UserChatRepoImpl{db: db}
}
