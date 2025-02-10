package domain

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type Chat struct {
	ID          uint64
	Name        string
	Type        ChatType
	Image       domain.Image
	LastMessage *Message
	CreatedBy   uint64
	Creator     *domain.User
	UserChats   []UserChat
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
