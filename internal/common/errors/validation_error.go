package errors

const ValidationErrorType = "ValidationError"

type ValidationError struct {
	*ErrorData
}

func NewValidationError(domain string, err error, data map[string]any, devDetails ...string) *ValidationError {
	return &ValidationError{
		ErrorData: NewErrorData(domain, ValidationErrorType, err, data, devDetails...),
	}
}
