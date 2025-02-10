package errors

import (
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/user/common"
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
