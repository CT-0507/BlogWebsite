package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/go-playground/validator/v10"
)

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	var hasUpper, hasLetter, hasDigit, hasSpecial bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
			hasLetter = true
		case unicode.IsLower(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLetter && hasDigit && hasSpecial
}

func ValidateStruct(lang string, s any) []messages.APIErrorResponse {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []messages.APIErrorResponse

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			errors = append(errors, messages.APIErrorResponse{
				Field:   e.Field(),
				Message: messages.MsgRequiredField.FormatLang(lang, e.Field()),
			})
		case "min":
			errors = append(errors, messages.APIErrorResponse{
				Field:   e.Field(),
				Message: messages.MsgMinLength.FormatLang(lang, e.Field(), e.Param()),
			})
		case "max":
			errors = append(errors, messages.APIErrorResponse{
				Field:   e.Field(),
				Message: messages.MsgMaxLength.FormatLang(lang, e.Field(), e.Param()),
			})
		default:
			errors = append(errors, messages.APIErrorResponse{
				Field:   e.Field(),
				Message: messages.MsgMaxLength.FormatLang(lang, e.Field(), e.Param()),
			})
		}

	}

	return errors
}
