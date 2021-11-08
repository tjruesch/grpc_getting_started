package server

import (
	"context"
	"log"
	"os"
	"unicode/utf8"

	"github.com/getsentry/sentry-go"
	protos "github.com/truesch/grpc_getting_started/protos/translation"
	"github.com/truesch/grpc_getting_started/vendors"
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

	var c vendors.Client
	switch i.GetVendor() {
	case protos.Vendors_DeepL:
		c = vendors.NewDeepLClient(os.Getenv("DEEPL_API_KEY"))
	default:
		c = vendors.NewGoogleClient(os.Getenv("GOOGLE_PROJECT_ID"))
	}

	log.Printf(
		`Translating Text "%s" (%s -> %s) with %s client.`,
		i.GetText(), i.GetSourceLang(), i.GetTargetLang(), i.GetVendor(),
	)

	resp, err := c.TranslateText([]string{i.GetText()}, i.GetSourceLang().String(), i.TargetLang.String())
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	tra := &protos.TranslationOutput{
		Text:        resp[0],
		SourceLang:  i.SourceLang,
		TargetLang:  i.TargetLang,
		BilledChars: int32(utf8.RuneCountInString(i.GetText())),
	}

	return tra, nil
}
