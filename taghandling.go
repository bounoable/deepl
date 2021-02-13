package deepl

const (
	// DefaultTagHandling is the default tag handling strategy: the translation
	// engine does not take tags into account.
	DefaultTagHandling = TagHandlingStrategy("default")
	// XMLTagHandling makes the API process XML input by extracting text out of
	// the structure, splitting it into individual sentences, translating them,
	// and placing them back into the XML structure.
	XMLTagHandling = TagHandlingStrategy("xml")
)

// TagHandlingStrategy is a `tag_handling` option.
type TagHandlingStrategy string

// Value returns the request value for f.
func (f TagHandlingStrategy) Value() string {
	if f == DefaultTagHandling {
		return ""
	}
	return string(f)
}

func (f TagHandlingStrategy) String() string {
	return string(f)
}
