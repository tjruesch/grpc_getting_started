package vendors

import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/getsentry/sentry-go"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

// Error Handling
type GoogleError struct {
	e string
}

func (m *GoogleError) Error() string {
	return m.e
}

func NewGoogleError(e string) *GoogleError {
	return &GoogleError{e: e}
}

// Client
type GoogleClient struct {
	ProjectID string
}

func NewGoogleClient(projectID string) *GoogleClient {
	return &GoogleClient{ProjectID: projectID}
}

// Translate a single string with the Google API
func (g *GoogleClient) TranslateText(
	text []string,
	sl string,
	tl string,
) ([]string, error) {

	ctx := context.Background()

	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return nil, fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", g.ProjectID),
		SourceLanguageCode: sl,
		TargetLanguageCode: tl,
		MimeType:           "text/plain", // Mime types: "text/plain", "text/html"
		Contents:           text,
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("TranslateText: %v", err)
	}

	// Display the translation for each input text provided
	var tr []string
	for _, translation := range resp.GetTranslations() {
		tr = append(tr, translation.GetTranslatedText())
	}

	return tr, nil
}

func (g *GoogleClient) TranslateFile(file io.ReadCloser, sl string, tl string) (io.ReadCloser, error) {
	return nil, NewGoogleError("Not Implemented")
}
