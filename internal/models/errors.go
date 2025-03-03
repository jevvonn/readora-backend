package models

import (
	"net/http"

	"github.com/jevvonn/readora-backend/internal/infra/errorpkg"
)

var (
	ErrEmailOrUsernameExists = errorpkg.NewError(
		"user with this username or email already exists", http.StatusBadRequest,
	)

	ErrEmailOrUsernameNotExists = errorpkg.NewError(
		"user with this username or email doesn't exists", http.StatusBadRequest,
	)

	ErrEmailNotExists = errorpkg.NewError(
		"user with this email doesn't exists", http.StatusBadRequest,
	)

	ErrEmailNotVerified = errorpkg.NewError(
		"email not verified", http.StatusBadRequest,
	)

	ErrEmailAlreadyVerified = errorpkg.NewError(
		"email already verified", http.StatusBadRequest,
	)

	ErrInvalidOTP = errorpkg.NewError(
		"invalid otp", http.StatusBadRequest,
	)

	ErrorInvalidEmailOrPassword = errorpkg.NewError(
		"invalid email/username or password", http.StatusBadRequest,
	)

	ErrOTPSent = errorpkg.NewError(
		"wait for 3 minutes before sending another OTP", http.StatusBadRequest,
	)
)
