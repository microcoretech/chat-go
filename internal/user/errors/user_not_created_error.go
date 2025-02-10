package errors

import (
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/user/common"
)

const UserNotCreatedErrorType = "UserNotCreatedError"

type UserNotCreatedError struct {
	*errors.ErrorData
}

func NewUserNotCreatedError() *UserNotCreatedError {
	return &UserNotCreatedError{
		ErrorData: errors.NewErrorData(common.UserDomain, UserNotCreatedErrorType, nil, nil),
	}
}
