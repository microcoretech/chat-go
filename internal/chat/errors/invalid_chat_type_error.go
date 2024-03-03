package errors

import (
	"chat/internal/chat/common"
	"chat/internal/common/errors"
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
