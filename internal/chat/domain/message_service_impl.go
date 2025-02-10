package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type MessageServiceImpl struct {
	messageRepo         MessageRepo
	userServiceContract UserServiceContract
}

func (s *MessageServiceImpl) fillMessage(ctx context.Context, message *Message) error {
	if message == nil {
		return nil
	}

	var userIDs []uint64

	userIDs = append(userIDs, message.CreatedBy)

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

	message.Creator = usersMap[message.CreatedBy]

	return nil
}

func (s *MessageServiceImpl) GetMessages(ctx context.Context, filter *MessageFilter) ([]Message, uint64, error) {
	count, err := s.messageRepo.GetMessagesCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return nil, 0, nil
	}

	messages, err := s.messageRepo.GetMessages(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	for index := range messages {
		if err := s.fillMessage(ctx, &messages[index]); err != nil {
			return nil, 0, err
		}
	}

	return messages, count, nil
}

func (s *MessageServiceImpl) CreateMessage(ctx context.Context, newMessage Message) (*Message, error) {
	message, err := s.messageRepo.CreateMessage(ctx, newMessage, nil)
	if err != nil {
		return nil, err
	}

	if err := s.fillMessage(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageServiceImpl) UpdateMessageStatus(
	ctx context.Context,
	messageIDs []uint64,
	messageStatus MessageStatus,
) error {
	if err := s.messageRepo.UpdateMessageStatus(
		ctx, messageIDs, messageStatus, nil); err != nil {
		return err
	}

	return nil
}

func NewMessageServiceImpl(
	messageRepo MessageRepo,
	userServiceContract UserServiceContract,
) *MessageServiceImpl {
	return &MessageServiceImpl{
		messageRepo:         messageRepo,
		userServiceContract: userServiceContract,
	}
}
