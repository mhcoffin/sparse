package sparse

import (
	"github.com/shoenig/test"
	"testing"
)

func TestExactly(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:    "hello world",
			parser:   Exactly(""),
			expected: &Tree{Runes: []rune{}},
		}, {
			input:    "hello world",
			parser:   Exactly("hello"),
			expected: &Tree{Runes: []rune("hello")},
		}, {
			input:    "hello world",
			parser:   Exactly("hello "),
			expected: &Tree{Runes: []rune("hello ")},
		}, {
			input:    "hello world",
			parser:   Exactly("hello world globe"),
			expected: nil,
		}, {
			input:  "foo",
			parser: Exactly("foo"), expected: &Tree{Runes: []rune("foo")},
		}, {
			input:  "",
			parser: Exactly("foo"), expected: nil,
		}, {
			input:    "",
			parser:   Exactly(""),
			expected: &Tree{Runes: []rune("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestDigit(t *testing.T) {
	testCases := []struct {
		input    string
		expected *Tree
	}{
		{"", nil},
		{"0abc", &Tree{Runes: []rune{'0'}}},
		{"1abc", &Tree{Runes: []rune{'1'}}},
		{"2abc", &Tree{Runes: []rune{'2'}}},
		{"3abc", &Tree{Runes: []rune{'3'}}},
		{"4abc", &Tree{Runes: []rune{'4'}}},
		{"5abc", &Tree{Runes: []rune{'5'}}},
		{"6abc", &Tree{Runes: []rune{'6'}}},
		{"7abc", &Tree{Runes: []rune{'7'}}},
		{"8abc", &Tree{Runes: []rune{'8'}}},
		{"9abc", &Tree{Runes: []rune{'9'}}},
		{"foo", nil},
		{"x", nil},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			matchLen := Digit([]rune(tt.input))
			test.Eq(t, tt.expected, matchLen)
		})
	}
}

func TestDigits(t *testing.T) {
	tests := []struct {
		input    string
		expected *Tree
	}{
		{input: "123", expected: &Tree{Runes: []rune("123")}},
		{input: "", expected: nil},
		{input: "foo", expected: nil},
		{input: "987foo", expected: &Tree{Runes: []rune("987")}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, Digits([]rune(tt.input)))
		})
	}
}

func TestAny(t *testing.T) {
	testCases := []struct {
		input    string
		expected *Tree
	}{
		{input: "", expected: nil},
		{input: "x", expected: &Tree{Runes: []rune{'x'}}},
		{input: "xyzzy", expected: &Tree{Runes: []rune{'x'}}},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, Any([]rune(tt.input)))
		})
	}
}

func TestAnyOf(t *testing.T) {
	testCases := []struct {
		input    string
		parser   Parser
		expected string
	}{
		{
			input:    "0234",
			parser:   OneOf("0123456789abcdef"),
			expected: "0",
		},
		{
			input:    "f",
			parser:   OneOf("0123456789abcdef"),
			expected: "f",
		},
		{
			input:    "g",
			parser:   OneOf("0123456789abcdef"),
			expected: "",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)).String())
		})
	}

}

func TestSeq(t *testing.T) {
	testCases := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:    "",
			parser:   Seq(),
			expected: &Tree{Runes: []rune{}, Children: []*Tree{}},
		}, {
			input:    "abcde",
			parser:   Seq(),
			expected: &Tree{Runes: []rune{}, Children: []*Tree{}},
		}, {
			input:  "3xcde",
			parser: Seq(Digit.Tagged("digits"), Exactly("x").Tagged("digits")),
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{Runes: []rune{'3'}, Tag: "digits"},
					{Runes: []rune{'x'}, Tag: "digits"},
				},
			},
		}, {
			input:  "3xcde",
			parser: Seq(Digit.Tagged("dig"), Exactly("x")),
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{Runes: []rune{'3'}, Tag: "dig"},
				},
			},
		}, {
			input:  "foo314159xcde",
			parser: Seq(Exactly("foo"), Digits.Tagged("digits"), Exactly("x")),
			expected: &Tree{
				Runes: []rune("foo314159x"),
				Children: []*Tree{
					{Runes: []rune("314159"), Tag: "digits"},
				},
			},
		}, {
			input:    "314159xcde",
			parser:   Seq(Exactly("foo"), Digits.Tagged("dig"), Exactly("x")),
			expected: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestOpt(t *testing.T) {
	testCases := []struct {
		input  string
		parser Parser
		result *Tree
	}{
		{
			input:  "abcdefghijklmnopqrstuvwxyz",
			parser: Optional(Digit),
			result: &Tree{Runes: []rune{}},
		}, {
			input:  "8abcdefghijklmnopqrstuvwxyz",
			parser: Optional(Digit),
			result: &Tree{Runes: []rune{'8'}},
		}, {
			input:  "88abcdefghijklmnopqrstuvwxyz",
			parser: Optional(Digit),
			result: &Tree{Runes: []rune{'8'}},
		}, {
			input:  "abcdefghijklmnopqrstuvwxyz",
			parser: Optional(Exactly("abc")),
			result: &Tree{Runes: []rune("abc")},
		}, {
			input:  "abcdefghijklmnopqrstuvwxyz",
			parser: Optional(Exactly("abd")),
			result: &Tree{Runes: []rune{}},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.result, tt.parser([]rune(tt.input)))
		})
	}
}

func TestZeroOrMore(t *testing.T) {
	tests := []struct {
		input  string
		parser Parser
		result *Tree
	}{
		{
			input:  "Foobar",
			parser: ZeroOrMore(Digit),
			result: &Tree{Runes: []rune{}},
		}, {
			input:  "123",
			parser: ZeroOrMore(Digit),
			result: &Tree{
				Runes: []rune("123"),
			},
		}, {
			input:  "1,2,3,",
			parser: ZeroOrMore(Digit.Tagged("digit"), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,3,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: "digit"},
					{Runes: []rune("2"), Tag: "digit"},
					{Runes: []rune{'3'}, Tag: "digit"},
				},
			},
		}, {
			input:  "1,2,3",
			parser: ZeroOrMore(Digit.Tagged("digit"), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: "digit"},
					{Runes: []rune("2"), Tag: "digit"},
				},
			},
		}, {
			input:  "1,2,3",
			parser: ZeroOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
			},
		}, {
			input:  "1,2,3",
			parser: ZeroOrMore(Digit, Optional(Exactly(","))),
			result: &Tree{
				Runes: []rune("1,2,3"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.result, tt.parser([]rune(tt.input)))
		})
	}
}

func TestOneOrMore(t *testing.T) {
	tests := []struct {
		input  string
		parser Parser
		result *Tree
	}{
		{
			input:  "Foobar",
			parser: OneOrMore(Digit),
			result: nil,
		}, {
			input:  "123",
			parser: OneOrMore(Digit),
			result: &Tree{
				Runes: []rune("123"),
			},
		}, {
			input:  "1,2,3,",
			parser: OneOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,3,"),
			},
		}, {
			input:  "1,2,3",
			parser: OneOrMore(Digit.Tagged("d"), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: "d"},
					{Runes: []rune("2"), Tag: "d"},
				},
			},
		}, {
			input:  "1,2,3",
			parser: OneOrMore(Digit.Tagged("d"), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: "d"},
					{Runes: []rune("2"), Tag: "d"},
				},
			},
		}, {
			input:  "1,2,3",
			parser: OneOrMore(Digit.Tagged("digit"), Optional(Exactly(","))),
			result: &Tree{
				Runes: []rune("1,2,3"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: "digit"},
					{Runes: []rune("2"), Tag: "digit"},
					{Runes: []rune("3"), Tag: "digit"},
				},
			},
		}, {
			input:  "(123)(0)(234x)",
			parser: ZeroOrMore(Exactly("("), Digits.Tagged("digit"), Exactly(")")),
			result: &Tree{
				Runes: []rune("(123)(0)"),
				Children: []*Tree{
					{Runes: []rune("123"), Tag: "digit"},
					{Runes: []rune("0"), Tag: "digit"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.result, tt.parser([]rune(tt.input)))
		})
	}
}

func TestFirstOf(t *testing.T) {
	testCases := []struct {
		input  string
		parser Parser
		result *Tree
	}{
		{
			input:  "abcdefghijklmnopqrstuvwxyz",
			parser: FirstOf(Digit, Letter),
			result: &Tree{Runes: []rune{'a'}},
		}, {
			input:  "abcdefghijklmnopqrstuvwxyz",
			parser: FirstOf(Letter, Digit),
			result: &Tree{Runes: []rune{'a'}},
		}, {
			input:  " abcdefghijklmnopqrstuvwxyz",
			parser: FirstOf(Letter, Digit),
			result: nil,
		}, {
			input:  " abcdefghijklmnopqrstuvwxyz",
			parser: FirstOf(Letter, Digit, Space),
			result: &Tree{Runes: []rune{' '}},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.result, tt.parser([]rune(tt.input)))
		})
	}
}

func TestLookingAt(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input: "if then else",
			parser: Seq(
				Exactly("if").Tagged("if"),
				WS,
				LookingAt(Exactly("then")),
			),
			expected: &Tree{
				Runes: []rune("if "),
				Children: []*Tree{
					{Runes: []rune("if"), Tag: "if"},
				},
			},
		},
		{
			input: "if then else",
			parser: Seq(
				WS,
				Exactly("if").Tagged("if"),
				WS,
				LookingAt(Exactly("th")),
				Exactly("then").Tagged("then"),
				WS,
			),
			expected: &Tree{
				Runes: []rune("if then "),
				Children: []*Tree{
					{Runes: []rune("if"), Tag: "if"},
					{Runes: []rune("then"), Tag: "then"},
				},
			},
		},
		{
			input: "if then else",
			parser: Seq(
				WS,
				Exactly("if"),
				WS,
				LookingAt(Exactly("other")),
				Exactly("then"),
				WS,
			),
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestNotLookingAt(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:    "abcdefg",
			parser:   Exactly("abc"),
			expected: &Tree{Runes: []rune{'a', 'b', 'c'}},
		}, {
			input:    "abcdefg",
			parser:   Not(Exactly("abc")),
			expected: nil,
		}, {
			input:    "abcdefg",
			parser:   Optional(Exactly("abc")),
			expected: &Tree{Runes: []rune{'a', 'b', 'c'}},
		}, {
			input:    "cdefg",
			parser:   Optional(Exactly("abc")),
			expected: &Tree{Runes: []rune{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
			if tt.expected == nil {
				test.Eq(t, &Tree{Runes: []rune{}}, Not(tt.parser)([]rune(tt.input)))
			} else {
				test.Eq(t, nil, Not(tt.parser)([]rune(tt.input)))
			}
		})
	}
}

func TestLetters(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:    "a stitch in time",
			parser:   Letters,
			expected: &Tree{Runes: []rune("a")},
		}, {
			input:    "stitch in time",
			parser:   Letters,
			expected: &Tree{Runes: []rune("stitch")},
		}, {
			input:  "a stitch in time",
			parser: ZeroOrMore(WS, Letters.Tagged("word")),
			expected: &Tree{
				Runes: []rune("a stitch in time"),
				Children: []*Tree{
					{Runes: []rune("a"), Tag: "word"},
					{Runes: []rune("stitch"), Tag: "word"},
					{Runes: []rune("in"), Tag: "word"},
					{Runes: []rune("time"), Tag: "word"},
				},
			},
		}, {
			input:  "a stitch in time",
			parser: OneOrMore(WS, Letters.Tagged("word")),
			expected: &Tree{
				Runes: []rune("a stitch in time"),
				Children: []*Tree{
					{Runes: []rune("a"), Tag: "word"},
					{Runes: []rune("stitch"), Tag: "word"},
					{Runes: []rune("in"), Tag: "word"},
					{Runes: []rune("time"), Tag: "word"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected *Tree
	}{
		{
			input:    "  ",
			expected: &Tree{Runes: []rune(" ")},
		}, {
			input:    "  x",
			expected: &Tree{Runes: []rune(" ")},
		}, {
			input:    "abc",
			expected: nil,
		}, {
			input:    "",
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, Space([]rune(tt.input)))
		})
	}
}

func TestSpaces(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:  " 	foo",
			parser: Spaces,
			expected: &Tree{
				Runes: []rune(" 	"),
			},
		}, {
			input:  "  		  foo",
			parser: Spaces,
			expected: &Tree{
				Runes: []rune("  		  "),
			},
		}, {
			input:    "foo",
			parser:   Spaces,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestIgnoreCase(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{
			input:    "lower case",
			parser:   IgnoreCase("lower case"),
			expected: &Tree{Runes: []rune("lower case")},
		}, {
			input:    "lOWeR case",
			parser:   IgnoreCase("lower case"),
			expected: &Tree{Runes: []rune("lOWeR case")},
		}, {
			input:    "lower case",
			parser:   IgnoreCase("Lower Case"),
			expected: &Tree{Runes: []rune("lower case")},
		}, {
			input:    "lower than that",
			parser:   IgnoreCase("Lower Case"),
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)))
		})
	}
}

func TestCombinations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		parser   Parser
		expected string
	}{
		{
			name:  "comment",
			input: "/* whatever goes in * here */",
			parser: Seq(
				Exactly("/*"),
				ZeroOrMore(Not(Exactly("*/")), Any),
				Exactly("*/"),
			),
			expected: "/* whatever goes in * here */",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser([]rune(tt.input)).String())
		})
	}
}

func TestFoo(t *testing.T) {
	m := ZeroOrMore(FirstOf(Letters.Tagged("word"), Digits.Tagged("number")))([]rune("first123second456"))
	test.Eq(t, []*Tree{
		{Tag: "word", Runes: []rune("first")},
		{Tag: "number", Runes: []rune("123")},
		{Tag: "word", Runes: []rune("second")},
		{Tag: "number", Runes: []rune("456")},
	}, m.Children)

	z := ZeroOrMore(FirstOf(Letters.Tagged("word"), Digits.Tagged("number")))([]rune(",first123"))
	test.Eq(t, &Tree{Runes: []rune{}}, z)

	n := ZeroOrMore(Exactly("number"), Digits.Tagged("number"))([]rune("number17number11number"))
	test.Eq(t, []*Tree{
		{Tag: "number", Runes: []rune("17")},
		{Tag: "number", Runes: []rune("11")},
	}, n.Children)
}

func TestZeroOrMoreOf(t *testing.T) {
	tests := []struct {
		parser   Parser
		input    string
		expected string
	}{
		{parser: ZeroOrMoreOf("abcde"), input: "", expected: ""},
		{parser: ZeroOrMoreOf("abcde"), input: "abghi", expected: "ab"},
		{parser: ZeroOrMoreOf("abcde"), input: "aaaaaaaaaaa", expected: "aaaaaaaaaaa"},
		{parser: ZeroOrMoreOf(""), input: "aaaaaaaaaaa", expected: ""},
		{parser: ZeroOrMoreOf("✍️∰"), input: "aaaaaaaaaaa", expected: ""},
		{parser: ZeroOrMoreOf("✍️∰a"), input: "aaaaaaaaaaa✍️ghi", expected: "aaaaaaaaaaa✍️"},
		{parser: ZeroOrMoreOf("✍️∰a"), input: "✍️∰✍️∰∰", expected: "✍️∰✍️∰∰"},
		{parser: ZeroOrMoreOf("✍️∰a"), input: "✍️∰✍️∰∰xyzzy", expected: "✍️∰✍️∰∰"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			test.Eq(t, []rune(tt.expected), tt.parser([]rune(tt.input)).Runes)
		})
	}
}
