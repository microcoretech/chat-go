package errors

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/common/errors"
)

const InvalidChatTypeErrorType = "InvalidChatTypeError"

type InvalidChatTypeError struct {
	*errors.ErrorData
}

func NewInvalidChatTypeError() *InvalidChatTypeError {
	return &InvalidChatTypeError{
		ErrorData: errors.NewErrorData(common.ChatDomain, InvalidChatTypeErrorType, nil, nil),
	}
}
