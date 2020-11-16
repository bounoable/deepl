package deepl

// Supported languages
const (
	Chinese    = Language("ZH")
	Dutch      = Language("NL")
	English    = Language("EN")
	French     = Language("FR")
	German     = Language("DE")
	Italian    = Language("IT")
	Japanese   = Language("JA")
	Polish     = Language("PL")
	Portuguese = Language("PT")
	Russian    = Language("RU")
	Spanish    = Language("ES")
)

// Language is a deepl language code.
type Language string
