package contract

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type UserService interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error)
}
