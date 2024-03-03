package errors

import (
	"chat/internal/common/errors"
	"chat/internal/user/common"
)

const UserAlreadyCreatedErrorType = "UserAlreadyCreatedError"

type UserAlreadyCreatedError struct {
	*errors.ErrorData
}

func NewUserAlreadyCreatedError(data map[string]any) *UserAlreadyCreatedError {
	return &UserAlreadyCreatedError{
		ErrorData: errors.NewErrorData(common.UserDomain, UserAlreadyCreatedErrorType, nil, data),
	}
}
