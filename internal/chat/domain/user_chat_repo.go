package domain

import (
	"chat/internal/common/repository"
	"context"
)

type UserChatRepo interface {
	CreateUserChats(ctx context.Context, userChats []UserChat, tx repository.Tx) error
	DeleteUserChats(ctx context.Context, userChats []UserChat, tx repository.Tx) error
}
