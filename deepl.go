package deepl

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	glossaryURL  string
}

// A ClientOption configures a Client.
type ClientOption func(*Client)

// A TranslateOption configures a translation request.
type TranslateOption func(url.Values)

// Error is a DeepL error.
type Error struct {
	// The HTTP error code, returned by the DeepL API.
	Code int

	Body []byte
}

// BaseURL returns a ClientOption that sets the base url for requests.
func BaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
		c.translateURL = fmt.Sprintf("%s/translate", c.baseURL)
		c.glossaryURL = fmt.Sprintf("%s/glossaries", c.baseURL)
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

// ShowBilledChars returns a TranslateOption that asks DeepL to return the
// number of billed characters.
func ShowBilledChars(show bool) TranslateOption {
	return func(vals url.Values) {
		vals.Set("show_billed_characters", boolString(show))
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

// GlossaryID returns a TranslateOption that sets the `glossary_id` DeepL
// option.
func GlossaryID(glossaryID string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("glossary_id", glossaryID)
	}
}

// Context returns a TranslateOption that sets the `context` DeepL
// option.
func Context(context string) TranslateOption {
	return func(vals url.Values) {
		vals.Set("context", context)
	}
}

// New returns a Client that uses authKey as the DeepL authentication key.
func New(authKey string, opts ...ClientOption) *Client {
	c := Client{
		authKey: authKey,
		client:  http.DefaultClient,
	}

	// default base url
	BaseURL(V2)(&c)

	for _, opt := range opts {
		opt(&c)
	}

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

	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)
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

func errorFromResp(r *http.Response) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	return Error{
		Code: r.StatusCode,
		Body: b,
	}
}

// CreateGlossary as per
// https://www.deepl.com/docs-api/managing-glossaries/creating-a-glossary/
func (c *Client) CreateGlossary(ctx context.Context, name string, sourceLang, targetLang Language, entries []GlossaryEntry) (*Glossary, error) {
	vals := make(url.Values)
	vals.Set("name", name)
	vals.Set("source_lang", string(sourceLang))
	vals.Set("target_lang", string(targetLang))
	vals.Set("entries_format", "tsv")
	entriesTSV := make([]string, 0, len(entries))
	for _, entry := range entries {
		entriesTSV = append(entriesTSV, entry.Source+"\t"+entry.Target)
	}
	vals.Set("entries", strings.Join(entriesTSV, "\n"))

	req, err := http.NewRequestWithContext(ctx, "POST", c.glossaryURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errorFromResp(resp)
	}

	var response Glossary
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode deepl response: %w", err)
	}

	return &response, nil
}

// ListGlossaries as per
// https://www.deepl.com/docs-api/managing-glossaries/listing-glossaries/
func (c *Client) ListGlossaries(ctx context.Context) ([]Glossary, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.glossaryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResp(resp)
	}

	var response struct {
		Glossaries []Glossary `json:"glossaries"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode deepl response: %w", err)
	}

	return response.Glossaries, nil
}

// ListGlossary as per
// https://www.deepl.com/docs-api/managing-glossaries/listing-glossary-information/
func (c *Client) ListGlossary(ctx context.Context, glossaryID string) (*Glossary, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.glossaryURL+"/"+glossaryID, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResp(resp)
	}

	var response Glossary
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode deepl response: %w", err)
	}

	return &response, nil
}

// ListGlossaryEntries as per
// https://www.deepl.com/docs-api/managing-glossaries/listing-entries-of-a-glossary/
func (c *Client) ListGlossaryEntries(ctx context.Context, glossaryID string) ([]GlossaryEntry, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.glossaryURL+"/"+glossaryID+"/entries", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)
	req.Header.Add("Accept", "text/tab-separated-values")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResp(resp)
	}

	var entries []GlossaryEntry
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 tab-separated values, got %q", line)
		}
		entries = append(entries, GlossaryEntry{
			Source: parts[0],
			Target: parts[1],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// DeleteGlossary as per
// https://www.deepl.com/docs-api/managing-glossaries/deleing-a-glossary/
func (c *Client) DeleteGlossary(ctx context.Context, glossaryID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.glossaryURL+"/"+glossaryID, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+c.authKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errorFromResp(resp)
	}
	return nil
}

func (err Error) Error() string {
	switch err.Code {
	case 456:
		return "Quota exceeded. The character limit has been reached."
	default:
		if len(err.Body) > 0 {
			return fmt.Sprintf("unexpected HTTP status %s (%s)",
				http.StatusText(err.Code),
				strings.TrimSpace(string(err.Body)))
		}
		return fmt.Sprintf("unexpected HTTP status %s",
			http.StatusText(err.Code))
	}
}

func boolString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
