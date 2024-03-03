package errors

import (
	"chat/internal/common/errors"
	"chat/internal/user/common"
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
