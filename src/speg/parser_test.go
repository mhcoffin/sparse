package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		tag      string
		expected *Tree
	}{
		{
			name:     "EmptyInput",
			input:    []rune(""),
			tag:      "",
			expected: nil,
		},
		{
			name:     "OneInput",
			input:    []rune("a"),
			tag:      "",
			expected: &Tree{Match: []rune("a")},
		},
		{
			name:  "TaggedOption",
			input: []rune("b"),
			tag:   "testTag",
			expected: &Tree{
				Match: []rune("b"),
				Tag:   "testTag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := Any().Tagged("tt.tag")
			context := NewContext()
			tree := matcher.Parse(tt.input, 0, context)
			if tt.expected != nil && tree != nil {
				tt.expected.Start = 0
				tt.expected.Tag = matcher.Tag()
			}
			test.Eq(t, tt.expected, tree)
			matcherCache, _ := context.cache[matcher.ID()]
			test.NotNil(t, matcherCache)
			cachedTree, ok := matcherCache[0]
			test.True(t, ok)
			test.Eq(t, tree, cachedTree)
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
			match := OneOrMoreLetters().Tagged(tt.name)
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
			matcher := OneOrMoreDigits().Tagged(tt.name)
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
	digitParser := OneOrMoreDigits().Tagged("digits")
	letterParser := OneOrMoreLetters().Tagged("letters")

	testCases := []struct {
		name           string
		parsers        Parser
		inputString    string
		start          int
		expectedResult *Tree
	}{
		{
			name:        "Single parser matches",
			parsers:     Sequence(digitParser),
			inputString: "123",
			expectedResult: &Tree{Match: []rune("123"), Children: []*Tree{
				{Match: []rune("123"), Tag: "digits"},
			}},
		},
		{
			name:        "Both parsers match",
			parsers:     Sequence(digitParser, letterParser).Tagged("test"),
			inputString: "123abc",
			expectedResult: &Tree{Match: []rune("123abc"), Tag: "test", Children: []*Tree{
				{Match: []rune("123"), Tag: "digits", Start: 0},
				{Match: []rune("abc"), Tag: "letters", Start: 3},
			}},
		},
		{
			name:           "Mismatch in parsers order",
			parsers:        Sequence(digitParser, letterParser),
			inputString:    "abc123",
			expectedResult: nil,
		},
		{
			name:           "No parsers provided",
			parsers:        Sequence(),
			inputString:    "123abc",
			expectedResult: &Tree{Match: []rune("")},
		},
		{
			name:        "start nonzero; with tag",
			parsers:     Sequence(digitParser, letterParser).Tagged("test"),
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
			name:   "single parser with tagged result succeeds",
			input:  "abcdef",
			parser: FirstOf(OneOrMoreLetters().Tagged("letters")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Children: []*Tree{
					{
						Start: 0,
						Match: []rune("abcdef"),
						Tag:   "letters",
					},
				},
			},
		},
		{
			name:   "single parser without tagged result succeeds",
			input:  "abcdef",
			parser: FirstOf(OneOrMoreLetters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:   "first parser succeeds",
			input:  "abcdef",
			parser: FirstOf(OneOrMoreLetters().Tagged("letters"), Digit()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Children: []*Tree{
					{
						Start: 0,
						Match: []rune("abcdef"),
						Tag:   "letters",
					},
				},
			},
		},
		{
			name:   "second parser succeeds with tag",
			input:  "abcdef",
			parser: FirstOf(Digit(), OneOrMoreLetters().Tagged("let")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Children: []*Tree{
					{
						Start: 0,
						Match: []rune("abcdef"),
						Tag:   "let",
					},
				},
			},
		},
		{
			name:   "second parser succeeds without tag",
			input:  "abcdef",
			parser: FirstOf(Digit(), OneOrMoreLetters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:     "no parser succeeds",
			input:    "+=abcdef",
			parser:   FirstOf(Digit(), OneOrMoreLetters()),
			expected: nil,
		},
		{
			name:     "empty fails",
			input:    "+=abcdef",
			parser:   FirstOf(),
			expected: nil,
		},
		{
			name:   "first parser succeeds start=3",
			input:  "+++abcdef",
			start:  3,
			parser: FirstOf(OneOrMoreLetters(), Digit()),
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
			parser: Opt(OneOrMoreLetters().Tagged("letters")),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
				Children: []*Tree{
					{
						Start: 0,
						Match: []rune("abcdef"),
						Tag:   "letters",
					},
				},
			},
		},
		{
			name:   "sub-parser succeeds without tag",
			input:  []rune("abcdef"),
			start:  0,
			parser: Opt(OneOrMoreLetters()),
			expected: &Tree{
				Start: 0,
				Match: []rune("abcdef"),
			},
		},
		{
			name:   "sub-parser fails",
			input:  []rune("abcdef"),
			parser: Opt(OneOrMoreDigits()),
			expected: &Tree{
				Start: 0,
				Match: []rune(""),
			},
		},
		{
			name:   "sub-parser succeeds pos=3",
			input:  []rune("abcdef"),
			start:  3,
			parser: Opt(OneOrMoreLetters()),
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


func TestLeftRecursion(t *testing.T) {
	var base = FirstOf(
		OneOrMoreLetters().Tagged("var"),
		OneOrMoreDigits().Tagged("num"),
	).Tagged("")

	iterm := Indirect()
	var term = FirstOf(
		Sequence(
			iterm.Tagged("lhs"),
			Exactly("*"),
			iterm.Tagged("rhs"),
		).Tagged("s"),
		base,
	).Tagged("term")
	iterm.Set(&term)

	tree := term.Parse([]rune("foobar*37"), 0, NewContext())
	test.Eq(t, "(first:term (seq:s [&term] `*` [&term]) (first letter+:var digit+:num))", term.String())
	test.Eq(t, "foobar*37", tree.Matched())
}
