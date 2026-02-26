package messages

import (
	"fmt"
)

type APIMessage string

const (
	ENGLISH  string = "en"
	JAPANESE string = "jp"
)

const (
	MsgRequiredField APIMessage = "REQUIRED_FIELD"
	MsgInvalidField  APIMessage = "INVALID_FIELD"
	MsgMinLength     APIMessage = "MIN_LENGTH"
	MsgMaxLength     APIMessage = "MAX_LENGTH"
	MsgFieldMisMatch APIMessage = "FIELD_MISMATCH"
)

var apiMessages = map[string]map[APIMessage]string{
	"en": {
		MsgRequiredField: "%s is required",
		MsgInvalidField:  "%s is invalid",
		MsgMinLength:     "%s must be at least %d characters",
		MsgMaxLength:     "%s cannot exceed %d characters",
		MsgFieldMisMatch: "%s is not matched",
	},
	"id": {
		MsgRequiredField: "%s wajib diisi",
	},
}

func (m APIMessage) FormatLang(lang string, args ...any) string {
	if l, ok := apiMessages[lang]; ok {
		if tpl, ok := l[m]; ok {
			return fmt.Sprintf(tpl, args...)
		}
	}
	return fmt.Sprintf(apiMessages["en"][m], args...)
}

type APIErrorResponse struct {
	Field   string
	Message string
}
