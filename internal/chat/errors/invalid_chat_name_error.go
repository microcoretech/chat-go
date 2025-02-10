package errors

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/common/errors"
)

const InvalidChatNameErrorType = "InvalidChatNameError"

type InvalidChatNameError struct {
	*errors.ErrorData
}

func NewInvalidChatNameError() *InvalidChatNameError {
	return &InvalidChatNameError{
		ErrorData: errors.NewErrorData(common.ChatDomain, InvalidChatNameErrorType, nil, nil),
	}
}
