package errors

import (
	"github.com/pkg/errors"
)

type (
	NotAllowedError struct {
		baseError
	}
	NotFoundError struct {
		baseError
	}
	NotValidError struct {
		baseError
	}
	StoreMismatchError struct {
		baseError
	}
	InternalError struct {
		baseError
	}
	UnknownError struct {
		baseError
	}
	baseError struct {
		error
		Msg string // Used to hide wrapped error when returning a response to client
	}
)

func NewNotAllowedError(msg string, err error) NotAllowedError {
	if err == nil {
		return NotAllowedError{
			baseError{errors.New(msg), msg},
		}
	}
	return NotAllowedError{
		baseError{errors.Wrap(err, msg), msg},
	}
}

func NewNotFoundError(msg string, err error) NotFoundError {
	if err == nil {
		return NotFoundError{
			baseError{errors.New(msg), msg},
		}
	}
	return NotFoundError{
		baseError{errors.Wrap(err, msg), msg},
	}
}

func NewNotValidError(msg string, err error) NotValidError {
	if err == nil {
		return NotValidError{
			baseError{errors.New(msg), msg},
		}
	}
	return NotValidError{
		baseError{errors.Wrap(err, msg), msg},
	}
}

func NewStoreMismatchError(msg string, err error) StoreMismatchError {
	if err == nil {
		return StoreMismatchError{
			baseError{error: errors.New(msg), Msg: msg},
		}
	}
	return StoreMismatchError{
		baseError{error: errors.Wrap(err, msg), Msg: msg},
	}
}

func NewInternalError(msg string, err error) InternalError {
	if err == nil {
		return InternalError{
			baseError{errors.New(msg), msg},
		}
	}
	return InternalError{
		baseError{errors.Wrap(err, msg), msg},
	}
}

func NewUnknownError(msg string, err error) UnknownError {
	if err == nil {
		return UnknownError{
			baseError{errors.New(msg), msg},
		}
	}
	return UnknownError{
		baseError{errors.Wrap(err, msg), msg},
	}
}
