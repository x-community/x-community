package handler

import (
	"github.com/micro/go-micro/v2/errors"
)

// Errors represents user service errors
var (
	errInvalidEmail                = errors.Error{Code: 1001, Detail: "Invalid email address"}
	errInvalidUsername             = errors.Error{Code: 1002, Detail: "Username required and must be less than 32 characters"}
	errInvalidPassword             = errors.Error{Code: 1003, Detail: "Password must have at least 8 characters"}
	errEmailAlreadyRegistered      = errors.Error{Code: 1004, Detail: "Email already registered"}
	errUsernameAlreadyRegistered   = errors.Error{Code: 1005, Detail: "Username already registered"}
	errIncorrectUsernameOrPassword = errors.Error{Code: 1006, Detail: "Username or password is incorrect"}
	errInvalidOldPassword          = errors.Error{Code: 1007, Detail: "Old password is incorrect"}
)

// InvalidAuthenticationCredential represents invalid authentication credential error
func (s *userService) InvalidAuthenticationCredential(err string) error {
	return &errors.Error{Id: s.id, Code: 22003, Detail: err}
}

// InternalServerError represents internal server error
func (s *userService) InternalServerError(err string) error {
	return errors.InternalServerError(s.id, err)
}

// NewError return a new micro service error
func (s *userService) NewError(err errors.Error) error {
	return &errors.Error{
		Id:     s.id,
		Code:   err.Code,
		Detail: err.Detail,
	}
}
