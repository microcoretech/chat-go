package errors

import (
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/user/common"
)

const UserNotFoundErrorType = "UserNotFoundError"

type UserNotFoundError struct {
	*errors.ErrorData
}

func NewUserNotFoundError(data map[string]any) *UserNotFoundError {
	return &UserNotFoundError{
		ErrorData: errors.NewErrorData(common.UserDomain, UserNotFoundErrorType, nil, data),
	}
}
