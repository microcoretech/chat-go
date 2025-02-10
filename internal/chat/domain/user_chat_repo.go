package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/repository"
)

type UserChatRepo interface {
	CreateUserChats(ctx context.Context, userChats []UserChat, tx repository.Tx) error
	DeleteUserChats(ctx context.Context, userChats []UserChat, tx repository.Tx) error
}
