package custom_errors

type BadRequestError struct {
	InternalError error
}

type ConflictError struct {
	InternalError error
}

func (c ConflictError) Error() string {
	return c.InternalError.Error()
}

func (b BadRequestError) Error() string {
	return b.InternalError.Error()
}
