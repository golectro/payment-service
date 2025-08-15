package utils

import (
	"errors"
	"fmt"
	"golectro-payment/internal/model"
	"strings"
)

func ParseMultilangError(err error) model.Message {
	msg := model.Message{}
	segments := strings.SplitSeq(err.Error(), "|")

	for segment := range segments {
		kv := strings.SplitN(strings.TrimSpace(segment), ":", 2)
		if len(kv) == 2 {
			lang := strings.ToLower(strings.TrimSpace(kv[0]))
			text := strings.TrimSpace(kv[1])
			msg[lang] = text
		}
	}

	return msg
}

func WrapMessageAsError(msg model.Message, err ...error) error {
	if len(err) > 0 && err[0] != nil {
		return err[0]
	}

	var segments []string
	for lang, text := range msg {
		segments = append(segments, fmt.Sprintf("%s: %s", lang, text))
	}

	return errors.New(strings.Join(segments, " | "))
}
