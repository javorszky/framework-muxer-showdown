package errors

import (
	"fmt"
)

// ApplicationError should probably be mapped to a 5xx response with custom messaging?
type ApplicationError struct {
	err error
}

func (a ApplicationError) Error() string {
	return fmt.Sprintf("application error: %s", a.err.Error())
}

// NotFoundError is mapped to a 404.
type NotFoundError struct {
	err error
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", n.err.Error())
}

// RequestError will map to one of the 4xx responses.
type RequestError struct {
	err error
}

func (r RequestError) Error() string {
	return fmt.Sprintf("request error: %s", r.err.Error())
}

// ShutdownError will result in the termination of the service, it should make it unhealthy.
type ShutdownError struct {
	err error
}

func (s ShutdownError) Error() string {
	return fmt.Sprintf("shutdown error: %s", s.err.Error())
}
