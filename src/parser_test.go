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
		{input: "hello world", parser: Exactly(""), expected: &Tree{Runes: []rune{}}},
		{input: "hello world", parser: Exactly("hello"), expected: &Tree{Runes: []rune("hello")}},
		{input: "hello world", parser: Exactly("hello "), expected: &Tree{Runes: []rune("hello ")}},
		{input: "hello world", parser: Exactly("hello world globe"), expected: nil},
		{input: "foo", parser: Exactly("foo"), expected: &Tree{Runes: []rune("foo")}},
		{input: "", parser: Exactly("foo"), expected: nil},
		{input: "", parser: Exactly(""), expected: &Tree{Runes: []rune("")}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

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
			parser: Seq(Digit.Tagged(9), Exactly("x").Tagged(10)),
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{Runes: []rune{'3'}, Tag: 9},
					{Runes: []rune{'x'}, Tag: 10},
				},
			},
		}, {
			input:  "3xcde",
			parser: Seq(Digit.Tagged(1), Exactly("x")),
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{Runes: []rune{'3'}, Tag: 1},
				},
			},
		}, {
			input:  "foo314159xcde",
			parser: Seq(Exactly("foo"), Digits.Tagged(1), Exactly("x")),
			expected: &Tree{
				Runes: []rune("foo314159x"),
				Children: []*Tree{
					{Runes: []rune("314159"), Tag: 1},
				},
			},
		}, {
			input:    "314159xcde",
			parser:   Seq(Exactly("foo"), Digits.Tagged(1), Exactly("x")),
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
			parser: ZeroOrMore(Digit.Tagged(1), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,3,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: 1},
					{Runes: []rune("2"), Tag: 1},
					{Runes: []rune{'3'}, Tag: 1},
				},
			},
		}, {
			input:  "1,2,3",
			parser: ZeroOrMore(Digit.Tagged(7), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: 7},
					{Runes: []rune("2"), Tag: 7},
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
			parser: OneOrMore(Digit.Tagged(2), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: 2},
					{Runes: []rune("2"), Tag: 2},
				},
			},
		}, {
			input:  "1,2,3",
			parser: OneOrMore(Digit.Tagged(1), Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: 1},
					{Runes: []rune("2"), Tag: 1},
				},
			},
		}, {
			input:  "1,2,3",
			parser: OneOrMore(Digit.Tagged(1), Optional(Exactly(","))),
			result: &Tree{
				Runes: []rune("1,2,3"),
				Children: []*Tree{
					{Runes: []rune("1"), Tag: 1},
					{Runes: []rune("2"), Tag: 1},
					{Runes: []rune("3"), Tag: 1},
				},
			},
		}, {
			input:  "(123)(0)(234x)",
			parser: ZeroOrMore(Exactly("("), Digits.Tagged(1), Exactly(")")),
			result: &Tree{
				Runes: []rune("(123)(0)"),
				Children: []*Tree{
					{Runes: []rune("123"), Tag: 1},
					{Runes: []rune("0"), Tag: 1},
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
				Exactly("if").Tagged(9),
				WS,
				LookingAt(Exactly("then")),
			),
			expected: &Tree{
				Runes: []rune("if "),
				Children: []*Tree{
					{Runes: []rune("if"), Tag: 9},
				},
			},
		},
		{
			input: "if then else",
			parser: Seq(
				WS,
				Exactly("if").Tagged(9),
				WS,
				LookingAt(Exactly("th")),
				Exactly("then").Tagged(11),
				WS,
			),
			expected: &Tree{
				Runes: []rune("if then "),
				Children: []*Tree{
					{Runes: []rune("if"), Tag: 9},
					{Runes: []rune("then"), Tag: 11},
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
			parser: ZeroOrMore(WS, Letters.Tagged(2)),
			expected: &Tree{
				Runes: []rune("a stitch in time"),
				Children: []*Tree{
					{Runes: []rune("a"), Tag: 2},
					{Runes: []rune("stitch"), Tag: 2},
					{Runes: []rune("in"), Tag: 2},
					{Runes: []rune("time"), Tag: 2},
				},
			},
		}, {
			input:  "a stitch in time",
			parser: OneOrMore(WS, Letters.Tagged(7)),
			expected: &Tree{
				Runes: []rune("a stitch in time"),
				Children: []*Tree{
					{Runes: []rune("a"), Tag: 7},
					{Runes: []rune("stitch"), Tag: 7},
					{Runes: []rune("in"), Tag: 7},
					{Runes: []rune("time"), Tag: 7},
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
