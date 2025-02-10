package errors

import (
	"mbobrovskyi/chat-go/internal/common/common"
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
