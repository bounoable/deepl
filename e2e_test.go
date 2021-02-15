package deepl_test

import (
	"context"
	"os"
	"testing"

	"github.com/bounoable/deepl"
	"github.com/stretchr/testify/assert"
)

func TestTranslate_withoutSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t))

	translated, sourceLang, err := client.Translate(
		context.Background(),
		"This is an example text.",
		deepl.German,
	)

	assert.Nil(t, err)
	assert.Equal(t, "Dies ist ein Beispieltext.", translated)
	assert.Equal(t, deepl.English, sourceLang)
}

func TestTranslate_withSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
		return
	}

	client := deepl.New(getAuthKey(t))

	_, sourceLang, err := client.Translate(
		context.Background(),
		"Voici un exemple de texte.",
		deepl.German,
		deepl.SourceLang(deepl.English),
	)

	assert.Nil(t, err)
	assert.Equal(t, deepl.English, sourceLang)

	// we don't validate the translated text, because the translation behaviour
	// for an invalid source language is not defined
}

func getAuthKey(t *testing.T) string {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		t.Fatal("Set the DEEPL_AUTH_KEY environment variable before running the integration tests.")
	}
	return authKey
}
