package naming

type NamingError struct {
	message string
}

func NewNamingError(message string) *NamingError {
	return &NamingError{
		message: message,
	}
}

func (e *NamingError) Error() string {
	return e.message
}
