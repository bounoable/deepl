package deepl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	httpi "github.com/bounoable/deepl/http"
)

const (
	// V2 is the base url for v2 of the deepl API.
	V2 = "https://api.deepl.com/v2"
)

// A Client is a deepl client.
type Client struct {
	client       httpi.Client
	authKey      string
	baseURL      string
	translateURL string
}

// A ClientOption configures a Client.
type ClientOption func(*Client)

// A TranslateOption configures a translation request.
type TranslateOption func(url.Values)

// Error is a DeepL error.
type Error struct {
	// The HTTP error code, returned by the deepl API.
	Code int
}

// BaseURL returns a ClientOption that sets the base url for requests.
func BaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// HTTPClient returns a ClientOption that specifies the http.Client that's used
// when making requests.
func HTTPClient(client httpi.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// SourceLang returns a ClientOption that specifies the source language of the
// input text. If SourceLang is not used, DeepL automatically figures out the
// source language.
func SourceLang(lang Language) TranslateOption {
	return func(vals url.Values) {
		vals.Set("source_lang", string(lang))
	}
}

// SplitSentences returns a TranslateOption that sets the `split_sentences`
// DeepL option.
func SplitSentences(split SplitSentence) TranslateOption {
	return func(vals url.Values) {
		vals.Set("split_sentences", split.Value())
	}
}

// PreserveFormatting returns a TranslateOption that sets the
// `preserve_formatting` DeepL option.
func PreserveFormatting(preserve bool) TranslateOption {
	return func(vals url.Values) {
		vals.Set("preserve_formatting", boolString(preserve))
	}
}

// Formality returns a TranslateOption that sets the `formality` DeepL option.
func Formality(formal Formal) TranslateOption {
	return func(vals url.Values) {
		vals.Set("formality", formal.Value())
	}
}

// TagHandling returns a TranslateOption that sets the `tag_handling` DeepL
// option.
func TagHandling(handling TagHandlingStrategy) TranslateOption {
	return func(vals url.Values) {
		vals.Set("tag_handling", handling.Value())
	}
}

// IgnoreTags returns a TranslateOption that sets the `ignore_tags` DeepL
// option.
func IgnoreTags(tags ...string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("ignore_tags", strings.Join(tags, ","))
	}
}

// New returns a Client that uses authKey as the DeepL authentication key.
func New(authKey string, opts ...ClientOption) *Client {
	c := Client{
		authKey: authKey,
		baseURL: V2,
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&c)
	}

	c.translateURL = fmt.Sprintf("%s/translate", c.baseURL)

	return &c
}

// HTTPClient returns the underlying http.Client.
func (c *Client) HTTPClient() httpi.Client {
	return c.client
}

// BaseURL returns the configured base url for requests.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// AuthKey returns the DeepL authentication key.
func (c *Client) AuthKey() string {
	return c.authKey
}

// Translate translates the provided text into the specified Language and
// returns the translated text and the detected source Language of the input
// text.
//
// When DeepL responds with an error, Translate returns an Error that contains
// the DeepL error code and message. Use errors.As to unwrap the returned error
// into an Error:
//
//	trans, sourceLang, err := c.Translate(context.TODO(), "Hello.", deepl.Japanese)
//	var deeplError deepl.Error
//	if errors.As(err, &deeplError) {
//		log.Println(fmt.Sprintf("DeepL error code %d: %s", deeplError.Code, deeplError))
//	}
func (c *Client) Translate(ctx context.Context, text string, targetLang Language, opts ...TranslateOption) (string, Language, error) {
	translations, err := c.TranslateMany(ctx, []string{text}, targetLang, opts...)
	if err != nil {
		return "", "", fmt.Errorf("translate many: %w", err)
	}

	if len(translations) == 0 {
		return "", "", errors.New("deepl responded with no translations")
	}

	return translations[0].Text, Language(translations[0].DetectedSourceLanguage), nil
}

// TranslateMany translates the provided texts into the specified Language and
// returns a Translation for every input text. The order of the translated texts
// is guaranteed to be the same as the order of the input texts.
//
// When DeepL responds with an error, TranslateMany returns an Error that
// contains the DeepL error code and message. Use errors.As to unwrap the
// returned error into an Error:
//
//	translations, err := c.TranslateMany(
//		context.TODO(),
//		[]string{"Hello", "World"},
//		deepl.Japanese,
//	)
//	var deeplError deepl.Error
//	if errors.As(err, &deeplError) {
//		log.Println(fmt.Sprintf("DeepL error code %d: %s", deeplError.Code, deeplError))
//	}
func (c *Client) TranslateMany(ctx context.Context, texts []string, targetLang Language, opts ...TranslateOption) ([]Translation, error) {
	vals := make(url.Values)
	vals.Set("auth_key", c.authKey)
	vals.Set("target_lang", string(targetLang))

	for _, text := range texts {
		vals.Add("text", text)
	}

	for _, opt := range opts {
		opt(vals)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.translateURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, Error{Code: resp.StatusCode}
	}

	var response translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode deepl response: %w", err)
	}

	return response.Translations, nil
}

func (err Error) Error() string {
	switch err.Code {
	case 456:
		return "Quota exceeded. The character limit has been reached."
	default:
		return http.StatusText(err.Code)
	}
}

func boolString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
