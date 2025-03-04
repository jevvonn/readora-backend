package errorpkg

import (
	"net/http"
)

var (
	ErrInternalServerError = NewError(
		"internal server error", http.StatusInternalServerError,
	)

	ErrEmailOrUsernameExists = NewError(
		"user with this username or email already exists", http.StatusBadRequest,
	)

	ErrEmailOrUsernameNotExists = NewError(
		"user with this username or email doesn't exists", http.StatusBadRequest,
	)

	ErrEmailNotExists = NewError(
		"user with this email doesn't exists", http.StatusBadRequest,
	)

	ErrEmailNotVerified = NewError(
		"email not verified", http.StatusBadRequest,
	)

	ErrEmailAlreadyVerified = NewError(
		"email already verified", http.StatusBadRequest,
	)

	ErrInvalidOTP = NewError(
		"invalid otp", http.StatusBadRequest,
	)

	ErrorInvalidEmailOrPassword = NewError(
		"invalid email/username or password", http.StatusBadRequest,
	)

	ErrOTPSent = NewError(
		"wait for 3 minutes before sending another OTP", http.StatusBadRequest,
	)
)
