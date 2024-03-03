package domain

import (
	"chat/internal/common/domain"
	"chat/internal/common/repository"
	"context"
)

type UserRepo interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUsersCount(ctx context.Context, filter *domain.UserFilter) (uint64, error)
	CreateUser(ctx context.Context, user domain.User, tx repository.Tx) (*domain.User, error)
}
