package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// applicationError should probably be mapped to a 5xx response with custom messaging?
type applicationError struct {
	err    error
	status int
}

func (a applicationError) Error() string {
	return fmt.Sprintf("application error: %s", a.err.Error())
}

// NewApplicationError returns a new initialised applicationError.
func NewApplicationError(err error, status int) error {
	return &applicationError{err: err, status: status}
}

// IsApplicationError checks whether the passed in error has an application error somewhere in its chain.
func IsApplicationError(err error) bool {
	var e *applicationError
	return errors.As(err, e)
}

// notFoundError is mapped to a 404.
type notFoundError struct {
	err    error
	status int
}

func (n notFoundError) Error() string {
	return fmt.Sprintf("not found: %s", n.err.Error())
}

// NewNotFoundError returns a new initialised notFoundError.
func NewNotFoundError(err error) error {
	return &notFoundError{err: err, status: http.StatusNotFound}
}

// IsNotFoundError checks whether the passed in error has a not found error somewhere in its chain.
func IsNotFoundError(err error) bool {
	var e *notFoundError
	return errors.As(err, e)
}

// requestError will map to one of the 4xx responses.
type requestError struct {
	err    error
	status int
}

func (r requestError) Error() string {
	return fmt.Sprintf("request error: %s", r.err.Error())
}

// NewRequestError returns an initialised requestError.
func NewRequestError(err error, status int) error {
	return &requestError{err: err, status: status}
}

// IsRequestError checks whether the passed in error has a request error somewhere in its chain.
func IsRequestError(err error) bool {
	var e *requestError
	return errors.As(err, e)
}

// shutdownError will result in the termination of the service, it should make it unhealthy.
type shutdownError struct {
	err    error
	status int
}

func (s shutdownError) Error() string {
	return fmt.Sprintf("shutdown error: %s", s.err.Error())
}

// NewShutdownError returns a new initialised shutdownError.
func NewShutdownError(err error, status int) error {
	return &shutdownError{err: err, status: status}
}

// IsShutdownError checks whether the passed in error has a shutdown error somewhere in its chain.
func IsShutdownError(err error) bool {
	var e *shutdownError
	return errors.As(err, e)
}
