package runtime_errors

type InternalServerError struct{
	Message string
}

func(e *InternalServerError) Error() string{
	return e.Message
}

type BadRequestError struct{
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

