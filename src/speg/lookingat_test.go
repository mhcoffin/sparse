package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestLookingAt(t *testing.T) {
	tests := []struct {
		name     string
		parser   Parser
		input    string
		start int
		expected string
	}{
		{"empty", LookingAt(Exactly("")), "", 0, `""`},
		{"empty", LookingAt(Exactly("foo")), "foobar", 0, `""`},
		{"empty", LookingAt(Exactly("bar")), "foobar", 0, `<nil>`},
		{"empty", LookingAt(Exactly("bar")), "foobar", 0, `<nil>`},
		{"empty", LookingAt(Exactly("bar")), "foobar", 3, `""`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, tc.parser.Parse([]rune(tc.input), tc.start, NewContext()).String())
		})
	}
}
