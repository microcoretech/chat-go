package domain

import "chat/internal/common/domain"

type UserCredentials struct {
	UserID   uint64
	Password string
	User     domain.User
}
