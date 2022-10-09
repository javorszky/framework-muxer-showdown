package errors

import (
	"errors"
	"fmt"
)

var (
	BaseAppError      = errors.New("some error from someplace")
	BaseNotFoundError = errors.New("not found the thing")
	BaseRequestError  = errors.New("hurr, bad request, yarr")
	BaseShutdownError = errors.New("unrecoverable error")
)

// ApplicationError should probably be mapped to a 5xx response with custom messaging?
type ApplicationError struct {
	err error
}

func (a ApplicationError) Error() string {
	return fmt.Sprintf("application error: %s", a.err.Error())
}

// NewApplicationError returns a new initialised applicationError.
func NewApplicationError(err error) *ApplicationError {
	return &ApplicationError{err: err}
}

// IsApplicationError checks whether the passed in error has an application error somewhere in its chain.
func IsApplicationError(err error) bool {
	var e *ApplicationError
	return errors.As(err, &e)
}

func GetApplicationError(err error) *ApplicationError {
	var e *ApplicationError
	if ok := errors.As(err, &e); !ok {
		return nil
	}
	return e
}

// NotFoundError is mapped to a 404.
type NotFoundError struct {
	err error
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", n.err.Error())
}

// NewNotFoundError returns a new initialised notFoundError.
func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err: err}
}

// IsNotFoundError checks whether the passed in error has a not found error somewhere in its chain.
func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

func GetNotFoundError(err error) *NotFoundError {
	var e *NotFoundError
	if ok := errors.As(err, &e); !ok {
		return nil
	}
	return e
}

// RequestError will map to one of the 4xx responses.
type RequestError struct {
	err error
}

func (r RequestError) Error() string {
	return fmt.Sprintf("request error: %s", r.err.Error())
}

// NewRequestError returns an initialised requestError.
func NewRequestError(err error) *RequestError {
	return &RequestError{err: err}
}

// IsRequestError checks whether the passed in error has a request error somewhere in its chain.
func IsRequestError(err error) bool {
	var e *RequestError
	return errors.As(err, &e)
}

func GetRequestError(err error) *RequestError {
	var e *RequestError
	if ok := errors.As(err, &e); !ok {
		return nil
	}
	return e
}

// ShutdownError will result in the termination of the service, it should make it unhealthy.
type ShutdownError struct {
	err error
}

func (s ShutdownError) Error() string {
	return fmt.Sprintf("shutdown error: %s", s.err.Error())
}

// NewShutdownError returns a new initialised shutdownError.
func NewShutdownError(err error) *ShutdownError {
	return &ShutdownError{err: err}
}

// IsShutdownError checks whether the passed in error has a shutdown error somewhere in its chain.
func IsShutdownError(err error) bool {
	var e *ShutdownError
	return errors.As(err, &e)
}

func GetShutdownError(err error) *ShutdownError {
	var e *ShutdownError
	if ok := errors.As(err, &e); !ok {
		return nil
	}
	return e
}
