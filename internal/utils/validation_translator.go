package utils

import (
	"golectro-payment/internal/constants"
	"golectro-payment/internal/model"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

func InitTranslator(v *validator.Validate) (enTrans, idTrans ut.Translator) {
	uni := ut.New(en.New(), en.New(), id.New())

	enTrans, _ = uni.GetTranslator("en")
	idTrans, _ = uni.GetTranslator("id")

	_ = en_translations.RegisterDefaultTranslations(v, enTrans)
	_ = id_translations.RegisterDefaultTranslations(v, idTrans)

	return
}

func TranslateValidationError(v *validator.Validate, err error) model.Message {
	enTrans, idTrans := InitTranslator(v)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return constants.FailedValidationOccurred
	}

	return model.Message{
		"id": joinMessages(validationErrors.Translate(idTrans)),
		"en": joinMessages(validationErrors.Translate(enTrans)),
	}
}

func joinMessages(msgs map[string]string) string {
	result := ""
	for _, m := range msgs {
		if result != "" {
			result += ", "
		}
		result += m
	}
	return result
}
