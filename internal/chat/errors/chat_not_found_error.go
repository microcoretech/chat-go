package errors

import (
	"chat/internal/chat/common"
	"chat/internal/common/errors"
)

const ChatNotFoundErrorType = "ChatNotFoundError"

type ChatNotFoundError struct {
	*errors.ErrorData
}

func NewChatNotFoundError(data map[string]any) *ChatNotFoundError {
	return &ChatNotFoundError{
		ErrorData: errors.NewErrorData(common.ChatDomain, ChatNotFoundErrorType, nil, data),
	}
}
