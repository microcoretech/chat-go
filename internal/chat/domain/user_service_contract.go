package domain

import (
	"chat/internal/common/domain"
	"context"
)

type UserServiceContract interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error)
}
