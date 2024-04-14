package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestMatcher_Star(t *testing.T) {
	tests := []struct {
		name     string
		parser   Parser
		input    string
		expected string
	}{
		{
			"empty", Letter().Star(), "", `""`,
		},

		{
			"one", Letter().Star(), "a", `"a"`,
		},
		{
			"several", Letter().Star(), "abcde", `"abcde"`,
		},
		{
			"several2", Letter().Star(), "abcde123", `"abcde"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser.Parse([]rune(tt.input), 0, NewContext()).String())
		})
	}
}
