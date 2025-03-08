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

	"golang.org/x/exp/maps"

	"chat-go/internal/chat/common"
	chaterrors "chat-go/internal/chat/errors"
	"chat-go/internal/common/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/repository"
)

type ChatServiceImpl struct {
	baseRepo            repository.BaseRepo
	chatRepo            ChatRepo
	userChatRepo        UserChatRepo
	userServiceContract UserServiceContract
}

func (s *ChatServiceImpl) fillChat(ctx context.Context, chat *Chat) error {
	if chat == nil {
		return nil
	}

	var userIDs []uint64

	userIDs = append(userIDs, chat.CreatedBy)

	if chat.LastMessage != nil {
		userIDs = append(userIDs, chat.LastMessage.CreatedBy)
	}

	for _, userChat := range chat.UserChats {
		userIDs = append(userIDs, userChat.UserID)
	}

	users, _, err := s.userServiceContract.GetUsers(ctx, &domain.UserFilter{
		IDs: userIDs,
	})
	if err != nil {
		return err
	}

	usersMap := make(map[uint64]*domain.User)
	for _, user := range users {
		usersMap[user.ID] = &user
	}

	chat.Creator = usersMap[chat.CreatedBy]

	if chat.LastMessage != nil {
		chat.LastMessage.Creator = usersMap[chat.LastMessage.CreatedBy]
	}

	for index := range chat.UserChats {
		chat.UserChats[index].User = usersMap[chat.UserChats[index].UserID]
	}

	return nil
}

func (s *ChatServiceImpl) GetChat(ctx context.Context, id uint64) (*Chat, error) {
	chat, err := s.chatRepo.GetChat(ctx, id)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.NewNotFoundError(common.ChatDomain)
	}

	if err := s.fillChat(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatServiceImpl) GetChats(ctx context.Context, filter *ChatFilter) ([]Chat, uint64, error) {
	count, err := s.chatRepo.GetChatsCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return nil, 0, nil
	}

	chats, err := s.chatRepo.GetChats(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for index := range chats {
		if err := s.fillChat(ctx, &chats[index]); err != nil {
			return nil, 0, err
		}
	}

	return chats, count, nil
}

func (s *ChatServiceImpl) CreateChat(ctx context.Context, chat Chat) (*Chat, error) {
	user := domain.UserFromContext(ctx)

	chat.CreatedBy = user.ID

	tx, err := s.baseRepo.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	uniqueUsers := make(map[uint64]struct{})
	uniqueUsers[chat.CreatedBy] = struct{}{}

	for _, userChat := range chat.UserChats {
		uniqueUsers[userChat.UserID] = struct{}{}
	}

	if chat.Type == DirectChatType {
		if len(uniqueUsers) != 2 {
			return nil, chaterrors.NewIncorrectUsersCountError()
		}
	} else {
		if len(chat.Name) == 0 {
			return nil, chaterrors.NewInvalidChatNameError()
		}
	}

	createdChat, err := s.chatRepo.CreateChat(ctx, chat, nil)
	if err != nil {
		return nil, err
	}

	userIDs := maps.Keys(uniqueUsers)
	userChats := make([]UserChat, len(userIDs))
	for index, id := range userIDs {
		userChats[index] = UserChat{
			UserID: id,
			ChatID: createdChat.ID,
		}
	}

	if err := s.userChatRepo.CreateUserChats(ctx, userChats, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetChat(ctx, createdChat.ID)
}

func (s *ChatServiceImpl) UpdateChat(ctx context.Context, chat Chat) (*Chat, error) {
	user := domain.UserFromContext(ctx)

	existingChat, err := s.GetChat(ctx, chat.ID)
	if err != nil {
		return nil, err
	}

	if existingChat == nil {
		return nil, errors.NewNotFoundError(common.ChatDomain)
	}

	if existingChat.CreatedBy != user.ID {
		return nil, errors.NewForbiddenError()
	}

	updatedChat, err := s.chatRepo.UpdateChat(ctx, chat)
	if err != nil {
		return nil, err
	}

	if updatedChat == nil {
		return nil, errors.NewNotFoundError(common.ChatDomain)
	}

	return updatedChat, nil
}

func (s *ChatServiceImpl) DeleteChat(ctx context.Context, id uint64) error {
	user := domain.UserFromContext(ctx)

	chat, err := s.chatRepo.GetChat(ctx, id)
	if err != nil {
		return err
	}

	if chat == nil {
		return errors.NewNotFoundError(common.ChatDomain)
	}

	if chat.CreatedBy != user.ID {
		return errors.NewForbiddenError()
	}

	err = s.chatRepo.DeleteChat(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func NewChatServiceImpl(
	baseRepo repository.BaseRepo,
	charRepo ChatRepo,
	userChatRepo UserChatRepo,
	userServiceContract UserServiceContract,
) *ChatServiceImpl {
	return &ChatServiceImpl{
		baseRepo:            baseRepo,
		chatRepo:            charRepo,
		userChatRepo:        userChatRepo,
		userServiceContract: userServiceContract,
	}
}
