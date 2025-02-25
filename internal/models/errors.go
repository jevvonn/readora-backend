package models

import "errors"

var (
	ErrEmailOrUsernameExists    = errors.New("user with this username or email already exists")
	ErrEmailOrUsernameNotExists = errors.New("user with this username or email doesn't exists")
	ErrEmailNotExists           = errors.New("user with this email doesn't exists")
	ErrEmailNotVerified         = errors.New("email not verified")
	ErrEmailAlreadyVerified     = errors.New("email already verified")
	ErrInvalidOTP               = errors.New("invalid otp")
)
