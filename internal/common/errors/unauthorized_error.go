package errors

import (
	"chat/internal/common/common"
)

const UnauthorizedErrorType = "UnauthorizedError"

type UnauthorizedError struct {
	*ErrorData
}

func NewUnauthorizedError(devDetails ...string) *UnauthorizedError {
	return &UnauthorizedError{
		ErrorData: NewErrorData(common.CommonDomain, UnauthorizedErrorType, nil, nil, devDetails...),
	}
}
