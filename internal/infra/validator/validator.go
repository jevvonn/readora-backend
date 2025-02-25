package validator

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationService interface {
	Validate(req interface{}) ([]map[string]string, error)
}

type Validator struct {
	vd         *validator.Validate
	translator ut.Translator
}

func NewValidator() ValidationService {
	vd := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(vd, trans)

	return &Validator{
		vd, trans,
	}
}

func (v *Validator) Validate(req interface{}) ([]map[string]string, error) {
	err := v.vd.Struct(req)
	if err != nil {
		errorMap := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			message := e.Translate(v.translator)
			errorMap = append(errorMap, map[string]string{
				"field":   e.Field(),
				"message": message,
			})
		}

		return errorMap, errors.New("validation errors")
	}

	return []map[string]string{}, nil
}
