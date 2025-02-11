// Copyright 2025 Mykhailo Bobrovskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"context"

	"chat-go/internal/common/domain"
	"chat-go/internal/common/repository"
	"chat-go/internal/user/errors"
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
