package deepl_test

import (
	"context"
	"os"
	"testing"

	"github.com/bounoable/deepl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslate_withoutSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	translated, sourceLang, err := client.Translate(
		context.Background(),
		"This is an example text.",
		deepl.German,
	)

	assert.Nil(t, err)
	assert.Equal(t, "Dies ist ein Beispieltext.", translated)
	assert.Equal(t, deepl.English, sourceLang)
}

func TestTranslate_showBilledCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	translations, err := client.TranslateMany(
		context.Background(),
		[]string{"This is an example text."},
		deepl.German,
		deepl.ShowBilledChars(true),
	)

	require.Nil(t, err)
	require.Len(t, translations, 1)
	assert.Equal(t, "Dies ist ein Beispieltext.", translations[0].Text)
	assert.Equal(t, deepl.English, deepl.Language(translations[0].DetectedSourceLanguage))
	assert.NotNil(t, translations[0].BilledCharacters)
	assert.True(t, translations[0].BilledCharacters > 0)
}

func TestTranslate_withSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	_, sourceLang, err := client.Translate(
		context.Background(),
		"Voici un exemple de texte.",
		deepl.German,
		deepl.SourceLang(deepl.English),
	)

	require.Nil(t, err)
	assert.Equal(t, deepl.English, sourceLang)

	// we don't validate the translated text, because the translation behaviour
	// for an invalid source language is not defined
}

func TestHTMLTagHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	res, _, err := client.Translate(
		context.Background(),
		`<p alt="This is a test.">This is a test.</p>`,
		deepl.German,
		deepl.TagHandling(deepl.HTMLTagHandling),
	)

	require.Nil(t, err)
	assert.Equal(t, `<p alt="This is a test.">Dies ist ein Test.</p>`, res)
}

func getOpts(opts ...deepl.ClientOption) []deepl.ClientOption {
	apiEndpoint := os.Getenv("DEEPL_API_ENDPOINT")
	ret := opts
	if apiEndpoint != "" {
		ret = append(ret, deepl.BaseURL(apiEndpoint))
	}
	return ret
}

func getAuthKey(t *testing.T) string {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		t.Fatal("Set the DEEPL_AUTH_KEY environment variable before running the integration tests.")
	}
	return authKey
}
