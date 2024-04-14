package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		parser   Parser
		expected string
	}{
		{
			name:     "EmptyInput",
			input:    []rune(""),
			parser:   Any(),
			expected: `<nil>`,
		},
		{
			name:     "OneInput",
			input:    []rune("a"),
			parser:   Any(),
			expected: `"a"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Eq(t, tt.expected, tt.parser.Parse(tt.input, 0, NewContext()).String())
		})
	}
}

func TestOneOrMoreLetters(t *testing.T) {
	tests := []struct {
		name  string
		input []rune
		want  *Tree
	}{
		{
			name:  "EmptyInput",
			input: []rune(""),
			want:  nil,
		},
		{
			name:  "SingleLetter",
			input: []rune("a"),
			want:  &Tree{Match: []rune("a"), Tag: "SingleLetter"},
		},
		{
			name:  "MultipleLetters",
			input: []rune("abc"),
			want:  &Tree{Match: []rune("abc"), Tag: "MultipleLetters"},
		},
		{
			name:  "LettersWithDigit",
			input: []rune("abc123"),
			want:  &Tree{Match: []rune("abc"), Tag: "LettersWithDigit"},
		},
		{
			name:  "StartWithDigit",
			input: []rune("1abc123"),
			want:  nil,
		},
		{
			name:  "NonAsciiLetters",
			input: []rune("абвгд"),
			want:  &Tree{Match: []rune("абвгд"), Tag: "NonAsciiLetters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := Letters().Tagged(tt.name)
			context := NewContext()
			got := match.Parse(tt.input, 0, context)
			test.Eq(t, tt.want, got)
		})
	}
}

func TestLetter(t *testing.T) {
	tests := []struct {
		name  string
		input []rune
		want  *Tree
	}{
		{
			name:  "SingleLetter",
			input: []rune("a"),
			want:  &Tree{Match: []rune("a")},
		},
		{
			name:  "MultiLetters",
			input: []rune("abc"),
			want:  &Tree{Match: []rune("a")},
		},
		{
			name:  "NoLetters",
			input: []rune("123"),
			want:  nil,
		},
		{
			name:  "EmptyInput",
			input: []rune(""),
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Letter().Parse(tt.input, 0, NewContext())
			test.Eq(t, tt.want, got)
		})
	}
}

func TestOneOrMoreDigits(t *testing.T) {
	tests := []struct {
		name  string
		input []rune
		want  *Tree
	}{
		{
			name:  "EmptyRunes",
			input: []rune(""),
			want:  nil,
		},
		{
			name:  "SingleDigit",
			input: []rune("1"),
			want: &Tree{
				Match: []rune("1"), Tag: "SingleDigit",
			},
		},
		{
			name:  "MultipleDigits",
			input: []rune("123"),
			want: &Tree{
				Match: []rune("123"), Tag: "MultipleDigits",
			},
		},
		{
			name:  "NonDigitStart",
			input: []rune("a123"),
			want:  nil,
		},
		{
			name:  "NonDigitEnd",
			input: []rune("123a"),
			want: &Tree{
				Match: []rune("123"), Tag: "NonDigitEnd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := Digits().Tagged(tt.name)
			got := matcher.Parse(tt.input, 0, NewContext())
			test.Eq(t, tt.want, got)
		})
	}
}

func TestDigit(t *testing.T) {
	type testCase struct {
		name  string
		input []rune
		want  *Tree
	}
	tests := []testCase{
		{
			name:  "EmptyInput",
			input: []rune(""),
			want:  nil,
		},
		{
			name:  "DigitFirst",
			input: []rune("2abc"),
			want:  &Tree{Match: []rune("2")},
		},
		{
			name:  "NonDigitFirst",
			input: []rune("a123"),
			want:  nil,
		},
		{
			name:  "MultipleDigitsFirst",
			input: []rune("123abc"),
			want:  &Tree{Match: []rune("1")},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mat := Digit()
			got := mat.Parse(tc.input, 0, NewContext())
			test.Eq(t, tc.want, got)
		})
	}
}

func TestSequenceParser_Parse(t *testing.T) {
	digitParser := Digits().Tagged("digits")
	letterParser := Letters().Tagged("letters")

	testCases := []struct {
		name           string
		parsers        Parser
		inputString    string
		start          int
		expectedResult *Tree
	}{
		{
			name:        "Single parser matches",
			parsers:     Seq(digitParser),
			inputString: "123",
			expectedResult: &Tree{Match: []rune("123"), Children: []*Tree{
				{Match: []rune("123"), Tag: "digits"},
			}},
		},
		{
			name:        "Both parsers match",
			parsers:     Seq(digitParser, letterParser).Tagged("test"),
			inputString: "123abc",
			expectedResult: &Tree{Match: []rune("123abc"), Tag: "test", Children: []*Tree{
				{Match: []rune("123"), Tag: "digits", Start: 0},
				{Match: []rune("abc"), Tag: "letters", Start: 3},
			}},
		},
		{
			name:           "Mismatch in parsers order",
			parsers:        Seq(digitParser, letterParser),
			inputString:    "abc123",
			expectedResult: nil,
		},
		{
			name:           "No parsers provided",
			parsers:        Seq(),
			inputString:    "123abc",
			expectedResult: &Tree{Match: []rune("")},
		},
		{
			name:        "start nonzero; with tag",
			parsers:     Seq(digitParser, letterParser).Tagged("test"),
			inputString: "+++123abc",
			start:       3,
			expectedResult: &Tree{Match: []rune("123abc"), Tag: "test", Start: 3, Children: []*Tree{
				{Match: []rune("123"), Tag: "digits", Start: 3},
				{Match: []rune("abc"), Tag: "letters", Start: 6},
			}},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seqParser := tc.parsers
			inputRunes := []rune(tc.inputString)
			result := seqParser.Parse(inputRunes, tc.start, NewContext())
			test.Eq(t, tc.expectedResult, result)
		})
	}
}

func TestFirstOf(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		start    int
		parser   Parser
		expected *Tree
	}{
		{
			name:   "single parser with newTaggedParser result succeeds",
			input:  "abcdef",
			parser: Or(Letters().Tagged("letters")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Tag:   "letters",
			},
		},
		{
			name:   "single parser without newTaggedParser result succeeds",
			input:  "abcdef",
			parser: Or(Letters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:   "first parser succeeds",
			input:  "abcdef",
			parser: Or(Letters().Tagged("letters"), Digit()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Tag:   "letters",
			},
		},
		{
			name:   "second parser succeeds with tag",
			input:  "abcdef",
			parser: Or(Digit(), Letters().Tagged("let")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Tag:   "let",
			},
		},
		{
			name:   "second parser succeeds without tag",
			input:  "abcdef",
			parser: Or(Digit(), Letters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:     "no parser succeeds",
			input:    "+=abcdef",
			parser:   Or(Digit(), Letters()),
			expected: nil,
		},
		{
			name:     "empty fails",
			input:    "+=abcdef",
			parser:   Or(),
			expected: nil,
		},
		{
			name:   "first parser succeeds start=3",
			input:  "+++abcdef",
			start:  3,
			parser: Or(Letters(), Digit()),
			expected: &Tree{
				Start: 3,
				Match: []rune("abcdef"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.parser.Parse([]rune(tc.input), tc.start, NewContext())
			test.Eq(t, tc.expected, result)
		})
	}
}

func TestOpt(t *testing.T) {
	testCases := []struct {
		name     string
		input    []rune
		start    int
		parser   Parser
		expected *Tree
	}{
		{
			name:   "subparser succeeds with tag",
			input:  []rune("abcdef"),
			start:  0,
			parser: Opt(Letters().Tagged("letters")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Tag:   "letters",
			},
		},
		{
			name:   "sub-parser succeeds without tag",
			input:  []rune("abcdef"),
			start:  0,
			parser: Opt(Letters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:   "sub-parser fails",
			input:  []rune("abcdef"),
			parser: Opt(Digits()),
			expected: &Tree{
				Start: 0,
			},
		},
		{
			name:   "sub-parser succeeds pos=3",
			input:  []rune("abcdef"),
			start:  3,
			parser: Opt(Letters()),
			expected: &Tree{
				Start: 3,
				Match: []rune("def"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.parser.Parse(tc.input, tc.start, NewContext())
			test.Eq(t, tc.expected, result)
		})
	}
}

func TestLeft(t *testing.T) {
	tests := []struct {
		name     string
		parser   Parser
		input    string
		expected string
	}{
		{
			name: "expression",
			parser: Left(
				Letters().Tagged("lhs"),
				Seq(
					Exactly("+").Tagged("op"),
					Letters().Tagged("rhs")),
			),
			input:    "a+b",
			expected: `((lhs "a") (op "+") (rhs "b"))`,
		},
		{
			name: "expression",
			parser: Left(
				Letters().Tagged("lhs"),
				Seq(
					Exactly("+").Tagged("op"),
					Letters().Tagged("rhs")),
			).Tagged("add"),
			input:    "a+b",
			expected: `(add (lhs "a") (op "+") (rhs "b"))`,
		},
		{
			name: "expression",
			parser: Left(
				Letters().Tagged("lhs"),
				Seq(Exactly("+").Tagged("op"), Letters().Tagged("rhs")),
			),
			input:    "a+b+c",
			expected: `(((lhs "a") (op "+") (rhs "b")) (op "+") (rhs "c"))`,
		},
		{
			name: "expression",
			parser: Left(
				Letters().Tagged("lhs"),
				Seq(
					Exactly("+").Tagged("op"),
					Letters().Tagged("rhs"),
				),
			).Tagged("add"),
			input:    "a+b+c",
			expected: `(add (add (lhs "a") (op "+") (rhs "b")) (op "+") (rhs "c"))`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test.Eq(t, tc.expected, tc.parser.Parse([]rune(tc.input), 0, NewContext()).String())
		})
	}
}

// func TestZeroOrMore(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		parser Parser
// 		input string
// 		expected string
// 	}{
// 		{"empty input", ZeroOrMore(Letter()), "", `""`},
// 		{"once", ZeroOrMore(Letter()), "a", `("a")`},
// 		{"twice", ZeroOrMore(Letter()), "aa", `("aa")`},
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			test.Eq(t, tc.expected, tc.parser.Parse([]rune(tc.input), 0, NewContext()).String())
// 		})
// 	}
// }
//
