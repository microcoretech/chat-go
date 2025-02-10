package domain

import "mbobrovskyi/chat-go/internal/common/domain"

type UserChat struct {
	UserID uint64 `json:"userId"`
	ChatID uint64 `json:"chatId"`

	User *domain.User
}
