package deepl

// Supported languages
const (
	Arabic             Language = "AR"
	Bulgarian          Language = "BG"
	ChineseSimplified  Language = "ZH-HANS"
	ChineseTraditional Language = "ZH-HANT"
	Czech              Language = "CS"
	Danish             Language = "DA"
	Dutch              Language = "NL"
	EnglishAmerican    Language = "EN-US"
	EnglishBritish     Language = "EN-GB"
	Estonian           Language = "ET"
	Finnish            Language = "FI"
	French             Language = "FR"
	German             Language = "DE"
	Greek              Language = "EL"
	Hungarian          Language = "HU"
	Italian            Language = "IT"
	Japanese           Language = "JA"
	Korean             Language = "KO"
	Latvian            Language = "LV"
	Lithuanian         Language = "LT"
	NorwegianBokmal    Language = "NB"
	Polish             Language = "PL"
	PortugueseBrazil   Language = "PT-BR"
	PortuguesePortugal Language = "PT-PT"
	Romanian           Language = "RO"
	Russian            Language = "RU"
	Slovak             Language = "SK"
	Slovenian          Language = "SL"
	Spanish            Language = "ES"
	Swedish            Language = "SV"
	Turkish            Language = "TR"
	Ukrainian          Language = "UK"
)

const (
	// English (unspecified).
	//
	// Deprecated: use EnglishAmerican or EnglishBritish instead.
	English Language = "EN"

	// Portuguese (unspecified).
	//
	// Deprecated: use PortugueseBrazil or PortuguesePortugal instead.
	Portuguese Language = "PT"

	// Chinese (unspecified).
	//
	// Deprecated: use ChineseSimplified or ChineseTraditional instead.
	Chinese Language = "ZH"
)

// Language is a deepl language code.
type Language string
