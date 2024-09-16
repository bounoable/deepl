package deepl

const (
	// SplitNone means no splitting at all, whole input is treated as one sentence.
	SplitNone SplitSentence = "0"
	// SplitDefault splits on interpunction and on newlines (default).
	SplitDefault SplitSentence = "1"
	// SplitNoNewlines  splits on interpunction only, ignoring newlines.
	SplitNoNewlines SplitSentence = "nonewlines"
)

// SplitSentence is a split_sentences option.
type SplitSentence string

// Value returns the request value for split.
func (split SplitSentence) Value() string {
	switch split {
	case SplitNone:
		return "0"
	case SplitDefault:
		return "1"
	case SplitNoNewlines:
		return "nonewlines"
	default:
		return "1"
	}
}

func (split SplitSentence) String() string {
	return split.Value()
}
