package deepl_test

import (
	"testing"

	"github.com/bounoable/deepl"
	"github.com/stretchr/testify/assert"
)

func TestSplitSentence_Value_String(t *testing.T) {
	tests := map[deepl.SplitSentence]string{
		deepl.SplitNone:       "0",
		deepl.SplitDefault:    "1",
		deepl.SplitNoNewlines: "nonewlines",
	}

	for split, v := range tests {
		assert.Equal(t, split.Value(), v)
		assert.Equal(t, split.String(), v)
	}
}
