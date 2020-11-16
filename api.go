package deepl

type translateResponse struct {
	Translations []Translation `json:"translations"`
}

// Translation is a translation result from deepl.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}
