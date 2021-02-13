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

// New returns a usable deepl client.
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

// A Client is a deepl client.
type Client struct {
	client       httpi.Client
	authKey      string
	baseURL      string
	translateURL string
}

// A ClientOption configures the deepl client.
type ClientOption func(*Client)

// BaseURL sets the base url that is used for requests.
func BaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// HTTPClient configures the *Client to use the given *http.Client for requests.
func HTTPClient(client httpi.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// HTTPClient returns the underlying *http.Client.
func (c *Client) HTTPClient() httpi.Client {
	return c.client
}

// BaseURL returns the configures base url for request.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// AuthKey returns the used authentication key.
func (c *Client) AuthKey() string {
	return c.authKey
}

// Translate translates the given text into the given targetLang.
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

// TranslateMany translates multiple texts into the given targetLang.
//
// Available options:
//	SourceLang(), SplitSentences(), PreserveFormatting(), Formality()
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

// A TranslateOption configures a translation.
type TranslateOption func(url.Values)

// SourceLang sets the source language of the text (`source_lang` option).
func SourceLang(lang Language) TranslateOption {
	return func(vals url.Values) {
		vals.Set("source_lang", string(lang))
	}
}

// SplitSentences sets the `split_sentences` option.
func SplitSentences(split SplitSentence) TranslateOption {
	return func(vals url.Values) {
		vals.Set("split_sentences", split.Value())
	}
}

// PreserveFormatting sets the `preserve_formatting` option.
func PreserveFormatting(preserve bool) TranslateOption {
	return func(vals url.Values) {
		vals.Set("preserve_formatting", boolString(preserve))
	}
}

// Formality sets the `formality` option.
func Formality(formal Formal) TranslateOption {
	return func(vals url.Values) {
		vals.Set("formality", formal.Value())
	}
}

// TagHandling sets the `tag_handling` option.
func TagHandling(handling TagHandlingStrategy) TranslateOption {
	return func(vals url.Values) {
		vals.Set("tag_handling", handling.Value())
	}
}

// IgnoreTags sets the `ignore_tags` option.
func IgnoreTags(tags ...string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("ignore_tags", strings.Join(tags, ","))
	}
}

func boolString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// Error is a deepl error.
type Error struct {
	// The HTTP error code, returned by the deepl API.
	Code int
}

func (err Error) Error() string {
	switch err.Code {
	case 456:
		return "Quota exceeded. The character limit has been reached."
	default:
		return http.StatusText(err.Code)
	}
}
