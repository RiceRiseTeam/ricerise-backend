package apperror

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
