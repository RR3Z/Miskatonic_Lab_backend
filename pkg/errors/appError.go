package errors

type AppError struct {
	Status  int
	Message string
	Err     error
}
