package errorpkg

import (
	"net/http"

	"github.com/jevvonn/readora-backend/internal/infra/validator"
)

var (
	ErrInternalServerError = NewError(
		"internal server error", http.StatusInternalServerError,
	)

	ErrBadRequest = NewError(
		"internal server error", http.StatusBadRequest,
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

	ErrValidationTimeFormat = func(field string) error {
		return validator.NewValidationErr([]validator.ErrorField{
			{
				Field:   field,
				Message: "Invalid date format, must be in ISO 8601 format",
			},
		}, "Invalid date format")
	}

	ErrValidationGenresArray = validator.NewValidationErr([]validator.ErrorField{
		{
			Field:   "genres",
			Message: "genres must be an array of string, e.g: ['Romance', 'Fiction']",
		},
	}, "genres must be an array of string, e.g: ['Romance', 'Fiction']")

	ErrValidationFileRequired = func(field string) error {
		return validator.NewValidationErr([]validator.ErrorField{
			{
				Field:   field,
				Message: "File is required",
			},
		}, "File is required")
	}

	ErrValidationFileMimeType = func(field string, types []string) error {
		mimeTypesJoined := ""
		for i, t := range types {
			if i == 0 {
				mimeTypesJoined = t
			} else {
				mimeTypesJoined = mimeTypesJoined + ", " + t
			}
		}

		msg := "File type is not allowed, allowed types: " + mimeTypesJoined

		return validator.NewValidationErr([]validator.ErrorField{
			{
				Field:   field,
				Message: msg,
			},
		}, msg)
	}
)
