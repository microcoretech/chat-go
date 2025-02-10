package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/repository"
)

type MessageRepo interface {
	GetMessages(ctx context.Context, filter *MessageFilter) ([]Message, error)
	GetMessagesCount(ctx context.Context, filter *MessageFilter) (uint64, error)
	CreateMessage(ctx context.Context, message Message, tx repository.Tx) (*Message, error)
	UpdateMessageStatus(
		ctx context.Context,
		messageIDs []uint64,
		messageStatus MessageStatus,
		tx repository.Tx,
	) error
}
