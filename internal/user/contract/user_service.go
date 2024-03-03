package contract

import (
	"chat/internal/common/domain"
	"context"
)

type UserService interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error)
}
