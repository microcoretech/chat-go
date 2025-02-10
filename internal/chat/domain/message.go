package domain

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type MessageStatus uint8

func (ms MessageStatus) ToUint8() uint8 {
	return uint8(ms)
}

const (
	DraftMessageStatus  MessageStatus = 1
	UnreadMessageStatus MessageStatus = 2
	ReadMessageStatus   MessageStatus = 3
)

type Message struct {
	ID        uint64
	Text      string
	Status    MessageStatus
	ChatID    uint64
	CreatedBy uint64
	Creator   *domain.User
	CreatedAt time.Time
	UpdatedAt time.Time
}
