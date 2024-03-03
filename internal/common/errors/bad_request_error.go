package errors

const BadRequestErrorType = "BadRequestError"

type BadRequestError struct {
	*ErrorData
}

func NewBadRequestError(domain string, err error, data map[string]any, devDetails ...string) *BadRequestError {
	return &BadRequestError{
		ErrorData: NewErrorData(domain, BadRequestErrorType, err, data, devDetails...),
	}
}
