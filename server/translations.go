package server

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/getsentry/sentry-go"
	protos "github.com/truesch/grpc_getting_started/protos/translation"
)

type Translation struct {
	protos.UnimplementedTranslationServer
}

func NewTranslation() *Translation {
	return &Translation{}
}

func (t *Translation) Translate(
	ctx context.Context,
	i *protos.TranslationInput,
) (*protos.TranslationOutput, error) {
	sentry.CaptureMessage(fmt.Sprintf(`Translating Text "%s" (%s -> %s)`, i.GetText(), i.GetSourceLang(), i.GetTargetLang()))

	tra := &protos.TranslationOutput{
		Text:        "Hello World",
		SourceLang:  i.GetSourceLang(),
		TargetLang:  i.GetTargetLang(),
		BilledChars: int32(utf8.RuneCountInString(i.GetText())),
	}

	return tra, nil
}
