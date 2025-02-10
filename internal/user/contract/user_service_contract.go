package contract

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type UserServiceContractImpl struct {
	userService UserService
}

func (c *UserServiceContractImpl) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	return c.userService.GetUser(ctx, id)
}

func (c *UserServiceContractImpl) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error) {
	return c.userService.GetUsers(ctx, filter)
}

func NewUserServiceContractImpl(userService UserService) *UserServiceContractImpl {
	return &UserServiceContractImpl{userService: userService}
}
