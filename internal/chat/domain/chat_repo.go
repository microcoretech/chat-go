package domain

import (
	"chat/internal/common/repository"
	"context"
)

type ChatRepo interface {
	GetChat(ctx context.Context, id uint64) (*Chat, error)
	GetChats(ctx context.Context, filter *ChatFilter) ([]Chat, error)
	GetChatsCount(ctx context.Context, filter *ChatFilter) (uint64, error)
	CreateChat(ctx context.Context, chat Chat, tx repository.Tx) (*Chat, error)
	UpdateChat(ctx context.Context, chat Chat) (*Chat, error)
	DeleteChat(ctx context.Context, id uint64) (bool, error)
}
