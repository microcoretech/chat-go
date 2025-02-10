package errors

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/common/errors"
)

const IncorrectUsersCountErrorType = "IncorrectUsersCountError"

type IncorrectUsersCountError struct {
	*errors.ErrorData
}

func NewIncorrectUsersCountError() *IncorrectUsersCountError {
	return &IncorrectUsersCountError{
		ErrorData: errors.NewErrorData(common.ChatDomain, IncorrectUsersCountErrorType, nil, nil),
	}
}
