package domain

import "mbobrovskyi/chat-go/internal/common/domain"

type UserCredentials struct {
	UserID   uint64
	Password string
	User     domain.User
}
