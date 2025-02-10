package errors

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/common/errors"
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
