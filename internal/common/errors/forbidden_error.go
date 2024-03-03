package errors

import (
	"chat/internal/common/common"
)

const ForbiddenErrorType = "ForbiddenError"

type ForbiddenError struct {
	*ErrorData
}

func NewForbiddenError() *ForbiddenError {
	return &ForbiddenError{
		ErrorData: NewErrorData(common.CommonDomain, ForbiddenErrorType, nil, nil),
	}
}
