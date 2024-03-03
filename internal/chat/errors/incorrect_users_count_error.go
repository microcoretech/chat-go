package errors

import (
	"chat/internal/chat/common"
	"chat/internal/common/errors"
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
