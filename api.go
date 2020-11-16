package deepl

type translateResponse struct {
	Translations []Translation `json:"translations"`
}

// Translation contains the translated text from deepl.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}
