package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestNot(t *testing.T) {
	tests := []struct {
		name string
		parser Parser
		input string
		expected string
	}{
		{"not letter", Not(Letter()), "1", `""`},
		{"not digit", Not(Digits()), "abc", `""`},
		{"not digit", Not(Digits()), "123", `<nil>`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, tc.parser.Parse([]rune(tc.input), 0, NewContext()).String())	
		})
	}
}

