package sparse

import (
	"fmt"
	"github.com/shoenig/test/must"
	"testing"
)

func TestExactly(t *testing.T) {
	testCases := []struct {
		input    string
		prefix   string
		expected *Tree
	}{
		{input: "hello world", prefix: "", expected: &Tree{Runes: []rune{}}},
		{input: "hello world", prefix: "hello", expected: &Tree{Runes: []rune("hello")}},
		{input: "hello world", prefix: "hello ", expected: &Tree{Runes: []rune("hello ")}},
		{input: "hello world", prefix: "hello world globe", expected: nil},
		{input: "foo", prefix: "foo", expected: &Tree{Runes: []rune("foo")}},
		{input: "", prefix: "foo", expected: nil},
	}
	for _, test := range testCases {
		t.Run(fmt.Sprintf(`match "%s" against "%s"`, test.prefix, test.input), func(t *testing.T) {
			input := []rune(test.input)
			result := Exactly(test.prefix)(input)
			must.Eq(t, test.expected, result)
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
	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			matchLen := Digit([]rune(test.input))
			must.Eq(t, test.expected, matchLen)
		})
	}
}

func TestDigits(t *testing.T) {
	tests := []struct {
		input string
		expected *Tree
	}{
		{
			input: "123",
			expected: &Tree{Runes: []rune("123")},
		},
		{
			input: "",
			expected: nil,
		},
		{
			input: "foo",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, Digits([]rune(test.input)))
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
	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, Any([]rune(test.input)))
		})
	}
}

func TestSeq(t *testing.T) {
	testCases := []struct {
		input    string ""
		matchers []Parser
		expected *Tree
	}{
		{"", []Parser{}, &Tree{Runes: []rune{}, Children: []*Tree{}}},
		{"abcde", []Parser{}, &Tree{Runes: []rune{}, Children: []*Tree{}}},
		{
			input: "3xcde",
			matchers: []Parser{
				Digit,
				Exactly("x"),
			},
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{
						Runes: []rune{'3'},
					},
					{
						Runes: []rune{'x'},
					},
				},
			},
		},
		{
			input: "3xcde",
			matchers: []Parser{
				Digit,
				Exactly("x").Elide(),
			},
			expected: &Tree{
				Runes: []rune("3x"),
				Children: []*Tree{
					{
						Runes: []rune{'3'},
					},
				},
			},
		},
		{
			input: "foo314159xcde",
			matchers: []Parser{
				Exactly("foo").Elide(),
				Digits,
				Exactly("x").Elide(),
			},
			expected: &Tree{
				Runes: []rune("foo314159x"),
				Children: []*Tree{
					{
						Runes: []rune("314159"),
					},
				},
			},
		},
		{
			input: "314159xcde",
			matchers: []Parser{
				Exactly("foo").Elide(),
				Digits,
				Exactly("x").Elide(),
			},
			expected: nil,
		},
	}
	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			result := Seq(test.matchers...)([]rune(test.input))
			must.Eq(t, test.expected, result)
		})
	}
}

func TestOpt(t *testing.T) {
	testCases := []struct {
		input   string
		matcher Parser
		result  *Tree
	}{
		{input: "abcdefghijklmnopqrstuvwxyz", matcher: Optional(Digit), result: &Tree{Runes: []rune{}}},
		{input: "8abcdefghijklmnopqrstuvwxyz", matcher: Optional(Digit), result: &Tree{Runes: []rune{'8'}}},
		{input: "88abcdefghijklmnopqrstuvwxyz", matcher: Optional(Digit), result: &Tree{Runes: []rune{'8'}}},
		{input: "abcdefghijklmnopqrstuvwxyz", matcher: Optional(Exactly("abc")), result: &Tree{Runes: []rune("abc")}},
		{input: "abcdefghijklmnopqrstuvwxyz", matcher: Optional(Exactly("abd")), result: &Tree{Runes: []rune{}}},
	}
	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.result, test.matcher([]rune(test.input)))
		})
	}
}

func TestZeroOrMore(t *testing.T) {
	tests := []struct {
		input   string
		matcher Parser
		result  *Tree
	}{
		{
			input:   "Foobar",
			matcher: ZeroOrMore(Digit),
			result:  &Tree{Runes: []rune{}},
		},
		{
			input:   "123",
			matcher: ZeroOrMore(Digit),
			result: &Tree{
				Runes: []rune("123"),
				Children: []*Tree{
					{Runes: []rune{'1'}},
					{Runes: []rune{'2'}},
					{Runes: []rune{'3'}},
				},
			},
		},
		{
			input:   "1,2,3,",
			matcher: ZeroOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,3,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune(",")},
					{Runes: []rune("2")},
					{Runes: []rune(",")},
					{Runes: []rune{'3'}},
					{Runes: []rune(",")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: ZeroOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune(",")},
					{Runes: []rune("2")},
					{Runes: []rune(",")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: ZeroOrMore(Digit, Exactly(",").Elide()),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune("2")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: ZeroOrMore(Digit, Optional(Exactly(",")).Elide()),
			result: &Tree{
				Runes: []rune("1,2,3"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune("2")},
					{Runes: []rune("3")},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.result, test.matcher([]rune(test.input)))
		})
	}
}

func TestOneOrMore(t *testing.T) {
	tests := []struct {
		input   string
		matcher Parser
		result  *Tree
	}{
		{
			input:   "Foobar",
			matcher: OneOrMore(Digit),
			result:  nil,
		},
		{
			input:   "123",
			matcher: OneOrMore(Digit),
			result: &Tree{
				Runes: []rune("123"),
				Children: []*Tree{
					{Runes: []rune{'1'}},
					{Runes: []rune{'2'}},
					{Runes: []rune{'3'}},
				},
			},
		},
		{
			input:   "1,2,3,",
			matcher: OneOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,3,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune(",")},
					{Runes: []rune("2")},
					{Runes: []rune(",")},
					{Runes: []rune{'3'}},
					{Runes: []rune(",")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: OneOrMore(Digit, Exactly(",")),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune(",")},
					{Runes: []rune("2")},
					{Runes: []rune(",")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: OneOrMore(Digit, Exactly(",").Elide()),
			result: &Tree{
				Runes: []rune("1,2,"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune("2")},
				},
			},
		},
		{
			input:   "1,2,3",
			matcher: OneOrMore(Digit, Optional(Exactly(",")).Elide()),
			result: &Tree{
				Runes: []rune("1,2,3"),
				Children: []*Tree{
					{Runes: []rune("1")},
					{Runes: []rune("2")},
					{Runes: []rune("3")},
				},
			},
		},
		{
			input:   "(123)(0)(234x)",
			matcher: ZeroOrMore(Exactly("(").Elide(), Digits, Exactly(")").Elide()),
			result: &Tree{
				Runes: []rune("(123)(0)"),
				Children: []*Tree{
					{Runes: []rune("123")},
					{Runes: []rune("0")},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.result, test.matcher([]rune(test.input)))
		})
	}
}

func TestFirstOf(t *testing.T) {
	testCases := []struct {
		input    string
		matchers []Parser
		result   *Tree
	}{
		{input: "abcdefghijklmnopqrstuvwxyz", matchers: []Parser{Digit, Letter}, result: &Tree{Runes: []rune{'a'}}},
		{input: "abcdefghijklmnopqrstuvwxyz", matchers: []Parser{Letter, Digit}, result: &Tree{Runes: []rune{'a'}}},
		{input: " abcdefghijklmnopqrstuvwxyz", matchers: []Parser{Letter, Digit}, result: nil},
		{input: " abcdefghijklmnopqrstuvwxyz", matchers: []Parser{Letter, Digit, Space}, result: &Tree{Runes: []rune{' '}}},
	}
	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			foo := FirstOf(test.matchers...)([]rune(test.input))
			must.Eq(t, test.result, foo)
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
				Exactly("if"),
				WS,
				LookingAt(Exactly("then")),
			),
			expected: &Tree{
				Runes: []rune("if "),
				Children: []*Tree{
					{Runes: []rune("if")},
				},
			},
		},
		{
			input: "if then else",
			parser: Seq(
				WS,
				Exactly("if"),
				WS,
				LookingAt(Exactly("th")),
				Exactly("then"),
				WS,
			),
			expected: &Tree{
				Runes: []rune("if then "),
				Children: []*Tree{
					{Runes: []rune("if")},
					{Runes: []rune("then")},
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
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, test.parser([]rune(test.input)))
		})
	}
}

func TestNotLookingAt(t *testing.T) {
	tests := []struct {
		input    string
		parser   Parser
		expected *Tree
	}{
		{input: "abcdefg", parser: Exactly("abc"), expected: &Tree{Runes: []rune{'a', 'b', 'c'}}},
		{input: "abcdefg", parser: Not(Exactly("abc")), expected: nil},
		{input: "abcdefg", parser: Optional(Exactly("abc")), expected: &Tree{Runes: []rune{'a', 'b', 'c'}}},
		{input: "cdefg", parser: Optional(Exactly("abc")), expected: &Tree{Runes: []rune{}}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, test.parser([]rune(test.input)))
			if test.expected == nil {
				must.Eq(t, &Tree{Runes: []rune{}}, Not(test.parser)([]rune(test.input)))
			} else {
				must.Eq(t, nil, Not(test.parser)([]rune(test.input)))
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
		},
		{
			input:    "stitch in time",
			parser:   Letters,
			expected: &Tree{Runes: []rune("stitch")},
		},
		{
			input:  "a stitch in time",
			parser: ZeroOrMore(WS, Letters),
			expected: &Tree{
				Runes: []rune("a stitch in time"),
				Children: []*Tree{
					{Runes: []rune("a")},
					{Runes: []rune("stitch")},
					{Runes: []rune("in")},
					{Runes: []rune("time")},
				},
			},
		},
		{
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
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, test.parser([]rune(test.input)))
		})
	}
}

func TestSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected *Tree
	}{
		{input: "  ", expected: &Tree{Runes: []rune(" ")}},
		{input: "  x", expected: &Tree{Runes: []rune(" ")}},
		{input: "abc", expected: nil},
		{input: "", expected: nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, Space([]rune(test.input)))
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
		},
		{
			input:  "  		  foo",
			parser: Spaces,
			expected: &Tree{
				Runes: []rune("  		  "),
			},
		},
		{
			input:    "foo",
			parser:   Spaces,
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			must.Eq(t, test.expected, test.parser([]rune(test.input)))
		})
	}
}

func TestIgnoreCase(t *testing.T) {
	tests := []struct {
		input    string
		target   string
		expected *Tree
	}{
		{
			input:    "lower case",
			target:   "lower case",
			expected: &Tree{Runes: []rune("lower case")},
		},
		{
			input:    "lOWeR case",
			target:   "lower case",
			expected: &Tree{Runes: []rune("lOWeR case")},
		},
		{
			input:    "lower case",
			target:   "Lower Case",
			expected: &Tree{Runes: []rune("lower case")},
		},
		{
			input:    "lower than that",
			target:   "Lower Case",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := IgnoreCase(test.target)([]rune(test.input))
			must.Eq(t, test.expected, got)
		})
	}
}
