package utils

import (
	"reflect"
	"strings"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/messages"
	"github.com/go-playground/validator/v10"
)

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
