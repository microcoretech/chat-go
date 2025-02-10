package errors

import (
	"mbobrovskyi/chat-go/internal/common/common"
)

const UndefinedErrorType = "UndefinedError"

type UndefinedError struct {
	*ErrorData
}

func NewUndefinedError(err error) *UndefinedError {
	return &UndefinedError{
		ErrorData: NewErrorData(common.CommonDomain, UndefinedErrorType, err, nil),
	}
}
