package http

import (
	"chat/internal/chat/domain"
	"context"
)

type MessageService interface {
	GetMessages(ctx context.Context, filter *domain.MessageFilter) ([]domain.Message, uint64, error)
}
