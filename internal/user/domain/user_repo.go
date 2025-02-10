package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/domain"
	"mbobrovskyi/chat-go/internal/common/repository"
)

type UserRepo interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUsersCount(ctx context.Context, filter *domain.UserFilter) (uint64, error)
	CreateUser(ctx context.Context, user domain.User, tx repository.Tx) (*domain.User, error)
}
