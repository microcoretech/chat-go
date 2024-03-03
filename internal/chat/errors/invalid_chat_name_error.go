package errors

import (
	"chat/internal/chat/common"
	"chat/internal/common/errors"
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
