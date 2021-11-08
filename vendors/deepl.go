package vendors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
)

// Error Handling
type DeepLError struct {
	e string
}

func (m *DeepLError) Error() string {
	return m.e
}

func NewDeepLError(e string) *DeepLError {
	return &DeepLError{e: e}
}

// Types
type deeplTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type deeplResponse struct {
	Translations []deeplTranslation `json:"translations"`
}

// Client
type DeepLClient struct {
	authKey string
}

func NewDeepLClient(authKey string) *DeepLClient {
	return &DeepLClient{authKey: authKey}
}

// TranslateText translates a given input text using the DeepL API
func (d *DeepLClient) TranslateText(
	text []string,
	sl string,
	tl string,
) ([]string, error) {

	BaseURL := os.Getenv("DEEPL_API_URL") + "/" + os.Getenv("DEEPL_API_VERSION") + "/translate"

	// define URL Parameters
	data := d.createRequest(text, sl, tl)

	body, err := d.callAPI(data, BaseURL)
	if err != nil {
		if strings.HasPrefix(err.Error(), "DeepL API Status not 200:") {
			// capture full exception incl. response in sentry
			sentry.CaptureException(err)
			// return only the DeepL error message to the client
			err := NewDeepLError(fmt.Sprintf("Vendor Error: %s", body))
			return nil, err
		}
		sentry.CaptureException(err)
		return nil, err
	}

	// parse translation response
	dr, err := d.parseResponse(body)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	// create list of strings from deeplResponse
	trans := []string{}
	for _, t := range dr.Translations {
		trans = append(trans, t.Text)
	}

	return trans, nil
}

func (d *DeepLClient) parseResponse(body []byte) (*deeplResponse, error) {
	dr := &deeplResponse{}
	err := json.Unmarshal(body, dr)

	return dr, err
}

func (d *DeepLClient) createRequest(text []string, sl string, tl string) *url.Values {
	data := url.Values{}

	for _, t := range text {
		data.Add("text", t)
	}
	data.Add("target_lang", tl)
	data.Add("source_lang", sl)

	return &data
}

// callAPI calls the DeepL API and returns the response body or any error that might occur
func (d *DeepLClient) callAPI(data *url.Values, url string) ([]byte, error) {
	// create new http client and prepare request
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", os.Getenv("DEEPL_API_KEY")))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check that status is 200 OK
	if resp.StatusCode != http.StatusOK {
		err := NewDeepLError(
			fmt.Sprintf("DeepL API Status not 200: %v, Response Body: %s", resp, body),
		)
		return body, err
	}

	return body, nil
}

func (d *DeepLClient) TranslateFile(file *os.File, sl string, tl string) (*os.File, error) {
	return nil, NewGoogleError("Not Implemented")
}
