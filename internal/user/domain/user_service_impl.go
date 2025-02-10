package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/domain"
	"mbobrovskyi/chat-go/internal/common/repository"
	"mbobrovskyi/chat-go/internal/user/errors"
)

type UserServiceImpl struct {
	baseRepo            repository.BaseRepo
	userRepo            UserRepo
	userCredentialsRepo UserCredentialsRepo
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewUserNotFoundError(map[string]any{"id": id})
	}

	return user, nil
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error) {
	count, err := s.userRepo.GetUsersCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return make([]domain.User, 0), count, nil
	}

	users, err := s.userRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserServiceImpl) GetUserCredentialsByUsername(ctx context.Context, username string) (*UserCredentials, error) {
	return s.userCredentialsRepo.GetUserCredentialsByUsername(ctx, username)
}

func NewUserServiceImpl(
	baseRepo repository.BaseRepo,
	userRepo UserRepo,
) *UserServiceImpl {
	return &UserServiceImpl{
		baseRepo: baseRepo,
		userRepo: userRepo,
	}
}
