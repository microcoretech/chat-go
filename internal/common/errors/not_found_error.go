package errors

const NotFoundErrorType = "NotFoundError"

type NotFoundError struct {
	*ErrorData
}

func NewNotFoundError(domain string) *NotFoundError {
	return &NotFoundError{
		ErrorData: NewErrorData(domain, NotFoundErrorType, nil, nil),
	}
}
