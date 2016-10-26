package models

import (
	"fmt"
)

func NewError(errType Error_Type, msg string) *Error {
	return &Error{
		Type:    errType,
		Message: msg,
	}
}

func ConvertError(err error) *Error {
	if err == nil {
		return nil
	}

	modelErr, ok := err.(*Error)
	if !ok {
		modelErr = NewError(Error_UnknownError, err.Error())
	}
	return modelErr
}

func (err *Error) ToError() error {
	if err == nil {
		return nil
	}
	return err
}

func (err *Error) Equal(other error) bool {
	if e, ok := other.(*Error); ok {
		if err == nil && e != nil {
			return false
		}
		return e.GetType() == err.GetType()
	}
	return false
}

func (err *Error) Error() string {
	return err.GetMessage()
}

var (
	ErrResourceNotFound = &Error{
		Type:    Error_ResourceNotFound,
		Message: "the requested resource could not be found",
	}

	ErrResourceExists = &Error{
		Type:    Error_ResourceExists,
		Message: "the requested resource already exists",
	}

	ErrResourceConflict = &Error{
		Type:    Error_ResourceConflict,
		Message: "the requested resource is in a conflicting state",
	}

	ErrDeadlock = &Error{
		Type:    Error_Deadlock,
		Message: "the request failed due to deadlock",
	}

	ErrBadRequest = &Error{
		Type:    Error_InvalidRequest,
		Message: "the request received is invalid",
	}

	ErrUnknownError = &Error{
		Type:    Error_UnknownError,
		Message: "the request failed for an unknown reason",
	}

	ErrDeserialize = &Error{
		Type:    Error_Deserialize,
		Message: "could not deserialize record",
	}

	ErrFailedToOpenEnvelope = &Error{
		Type:    Error_FailedToOpenEnvelope,
		Message: "could not open envelope",
	}

	ErrGUIDGeneration = &Error{
		Type:    Error_GUIDGeneration,
		Message: "cannot generate random guid",
	}
)

type ErrInvalidField struct {
	Field string
}

func (err ErrInvalidField) Error() string {
	return "Invalid field: " + err.Field
}

type ErrInvalidModification struct {
	InvalidField string
}

func (err ErrInvalidModification) Error() string {
	return "attempt to make invalid change to field: " + err.InvalidField
}

func NewTaskTransitionError(from, to State) *Error {
	return &Error{
		Type:    Error_InvalidStateTransition,
		Message: fmt.Sprintf("Cannot transition from %v to %v", from, to),
	}
}

func NewUnrecoverableError(err error) *Error {
	return &Error{
		Type:    Error_Unrecoverable,
		Message: fmt.Sprint("Unrecoverable Error: ", err),
	}
}
