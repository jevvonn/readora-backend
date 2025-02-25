package models

import "errors"

var (
	ErrEmailOrUsernameExists = errors.New("user with this username or email already exists")
	ErrEmailNotVerified      = errors.New("email not verified")
)
