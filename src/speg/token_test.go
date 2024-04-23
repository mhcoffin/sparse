package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestTokenParser(t *testing.T) {
	tests := []struct {
		name string
		parser Parser
		input string
		expected string
	}{
		{"letters", Token(Letters()), "abc def", `("abc")`},
		{"letters", Token(Letters()), "  abc def", `("abc")`},
		{"letters", Token(Digits()), "  123 abc def", `("123")`},
		{"var", Token(Seq(Letter(), Star(Or(Letter(), Digit())))), "  xyz123", `("xyz123")`},
		{"var", Token(Seq(Letter(), Star(Or(Letter(), Digit())))), "  xyz 123", `("xyz")`},
		{"or", Token(Or(Exactly("+"), Exactly("-"))), "+", `("+")`},
		{"or", Token(Or(Exactly("+"), Exactly("-"))), "-", `("-")`},
		{"or", Token(Or(Exactly("+"), Exactly("-"))), "  - ", `("-")`},
		{"sum", Seq(Token(Letters()), Token(Exactly("+")), Token(Letters())), "abc + def", `(("abc") ("+") ("def"))`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, tc.parser.Parse([]rune(tc.input), 0, NewContext()).String())
		})
	}
}

