package naming

import "fmt"

type NamingError struct {
	message string
}

func NewNamingError(message string) *NamingError {
	return &NamingError{
		message: fmt.Sprintf("[Naming] %s", message),
	}
}

func (e *NamingError) Error() string {
	return e.message
}
