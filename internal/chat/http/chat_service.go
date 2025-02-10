package http

import (
	"context"

	"mbobrovskyi/chat-go/internal/chat/domain"
)

type ChatService interface {
	GetChat(ctx context.Context, id uint64) (*domain.Chat, error)
	GetChats(ctx context.Context, filter *domain.ChatFilter) ([]domain.Chat, uint64, error)
	CreateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error)
	UpdateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error)
	DeleteChat(ctx context.Context, id uint64) error
}
