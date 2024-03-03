package errors

const DatabaseErrorType = "DatabaseError"

type DatabaseError struct {
	*ErrorData
}

func NewDatabaseError(domain string, err error, devDetails ...string) *DatabaseError {
	return &DatabaseError{
		ErrorData: NewErrorData(domain, DatabaseErrorType, err, nil, devDetails...),
	}
}
