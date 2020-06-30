package handler

import (
	"github.com/micro/go-micro/v2/errors"
)

var (
	errInvalidEmailAddress = errors.Error{Code: 1101, Detail: "Invalid email address"}
	errInvalidEmailContent = errors.Error{Code: 1102, Detail: "Invalid email content"}
)

// InternalServerError represents internal server error
func (s *service) InternalServerError(err string) error {
	return errors.InternalServerError(s.id, err)
}

// NewError return a new micro service error
func (s *service) NewError(err errors.Error) error {
	return &errors.Error{
		Id:     s.id,
		Code:   err.Code,
		Detail: err.Detail,
	}
}
