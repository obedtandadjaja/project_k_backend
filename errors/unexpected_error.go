package errors

type UnexpectedError struct {
	Message string
	Err     error
}

func (error UnexpectedError) Error() string {
	if error.Message != "" {
		return error.Message
	}

	return error.Err.Error()
}
