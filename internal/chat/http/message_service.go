package http

import (
	"context"

	"mbobrovskyi/chat-go/internal/chat/domain"
)

type MessageService interface {
	GetMessages(ctx context.Context, filter *domain.MessageFilter) ([]domain.Message, uint64, error)
}
