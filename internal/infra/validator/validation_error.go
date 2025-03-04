package validator

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors  []ErrorField `json:"errors"`
	Message string       `json:"message"`
}

func NewValidationErr(errs []ErrorField, msg string) *ValidationError {
	return &ValidationError{
		errs, msg,
	}
}

func (v *ValidationError) Error() string {
	return v.Message
}
