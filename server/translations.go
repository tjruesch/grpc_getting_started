package server

import (
	"context"
	"fmt"
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
		c = vendors.NewGoogleClient(os.Getenv("GOOGLE_PROJECT_ID"))
	case protos.Vendors_GoogleTranslate:
		c = vendors.NewGoogleClient(os.Getenv("GOOGLE_PROJECT_ID"))
	case protos.Vendors_MMT:
		c = vendors.NewGoogleClient(os.Getenv("GOOGLE_PROJECT_ID"))
	}

	sentry.CaptureMessage(
		fmt.Sprintf(`Translating Text "%s" (%s -> %s) with %s client.`,
			i.GetText(), i.GetSourceLang(), i.GetTargetLang(), i.GetVendor(),
		),
	)

	resp, err := c.TranslateText(i.GetText(), i.GetSourceLang().String(), i.TargetLang.String())
	if err != nil {
		sentry.CaptureException(err)
	}

	tra := &protos.TranslationOutput{
		Text:        resp[0],
		SourceLang:  i.GetSourceLang(),
		TargetLang:  i.GetTargetLang(),
		BilledChars: int32(utf8.RuneCountInString(i.GetText())),
	}

	return tra, nil
}
