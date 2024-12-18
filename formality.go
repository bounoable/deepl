package deepl

const (
	// DefaultFormal is the default formality.
	DefaultFormal Formal = "default"
	// LessFormal means the text is written in a less formal / more informal language.
	LessFormal Formal = "less"
	// MoreFormal means the text is written in a more formal language.
	MoreFormal Formal = "more"
)

// Formal is a formality option.
type Formal string

// Value returns the request value for f.
func (f Formal) Value() string {
	return string(f)
}

// String returns the formality level as a [string].
func (f Formal) String() string {
	return f.Value()
}
