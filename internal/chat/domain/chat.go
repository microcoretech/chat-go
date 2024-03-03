package domain

import (
	"chat/internal/common/domain"
	"time"
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
