package speg

import (
	"github.com/shoenig/test"
	"testing"
)

func TestAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		tag string
		expected *Tree
	}{
		{
			name:     "EmptyInput",
			input:    []rune(""),
			tag: "",
			expected: nil,
		},
		{
			name:     "OneInput",
			input:    []rune("a"),
			tag: "",
			expected: &Tree{MatchedRunes: []rune("a")},
		},
		{
			name: "WithTagOption",
			input: []rune("b"),
			tag: "testTag",
			expected: &Tree{
				MatchedRunes: []rune("b"),
				Tag:          "testTag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := Any().WithTag("tt.tag")
			cache := make(Cache)
			tree := matcher.Parse(tt.input, 0, cache)
			if tt.expected != nil && tree != nil {
				tt.expected.MatchStartIndex = 0
				tt.expected.Tag = matcher.tag
			}
			test.Eq(t, tt.expected, tree)
			matcherCache, _ := cache[matcher.id]
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
			want:  &Tree{MatchedRunes: []rune("a"), Tag: "SingleLetter"},
		},
		{
			name:  "MultipleLetters",
			input: []rune("abc"),
			want:  &Tree{MatchedRunes: []rune("abc"), Tag: "MultipleLetters"},
		},
		{
			name:  "LettersWithDigit",
			input: []rune("abc123"),
			want:  &Tree{MatchedRunes: []rune("abc"), Tag: "LettersWithDigit"},
		},
		{
			name:  "StartWithDigit",
			input: []rune("1abc123"),
			want:  nil,
		},
		{
			name:  "NonAsciiLetters",
			input: []rune("абвгд"),
			want:  &Tree{MatchedRunes: []rune("абвгд"), Tag: "NonAsciiLetters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := OneOrMoreLetters().WithTag(tt.name)
			cache := make(map[ID]map[int]*Tree)
			got := match.Parse(tt.input, 0, cache)
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
			want:  &Tree{MatchedRunes: []rune("a")},
		},
		{
			name:  "MultiLetters",
			input: []rune("abc"),
			want:  &Tree{MatchedRunes: []rune("a")},
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
			got := Letter().Parse(tt.input, 0, make(Cache))
			test.Eq(t, tt.want, got)
		})
	}
}

func TestOneOrMoreDigits(t *testing.T) {
	tests := []struct {
		name string
		input []rune
		want *Tree
	}{
		{
			name: "EmptyRunes",
			input: []rune(""),
			want: nil,
		},
		{
			name: "SingleDigit",
			input: []rune("1"),
			want: &Tree{
				MatchedRunes: []rune("1"), Tag: "SingleDigit",
			},
		},
		{
			name: "MultipleDigits",
			input: []rune("123"),
			want: &Tree{
				MatchedRunes: []rune("123"), Tag: "MultipleDigits",
			},
		},
		{
			name: "NonDigitStart",
			input: []rune("a123"),
			want: nil,
		},
		{
			name: "NonDigitEnd",
			input: []rune("123a"),
			want: &Tree{
				MatchedRunes: []rune("123"), Tag: "NonDigitEnd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := OneOrMoreDigits().WithTag(tt.name)
			got := matcher.Parse(tt.input, 0, NewCache())
			test.Eq(t, tt.want, got)
		})
	}
}


func TestDigit(t *testing.T) {
	type testCase struct {
		name string
		input []rune
		want *Tree
	}
	tests := []testCase{
		{
			name: "EmptyInput",
			input: []rune(""),
			want:  nil,
		},
		{
			name: "DigitFirst",
			input: []rune("2abc"),
			want:  &Tree{MatchedRunes: []rune("2")},
		},
		{
			name: "NonDigitFirst",
			input: []rune("a123"),
			want:  nil,
		},
		{
			name: "MultipleDigitsFirst",
			input: []rune("123abc"),
			want:  &Tree{MatchedRunes: []rune("1")},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mat := Digit()
			got := mat.Parse(tc.input, 0, NewCache())
			test.Eq(t, tc.want, got)
		})
	}
}

func TestSequenceParser_Parse(t *testing.T) {
	digitParser := OneOrMoreDigits().WithTag("digits")
	letterParser := OneOrMoreLetters().WithTag("letters")

	testCases := []struct {
		name           string
		parsers        Parser
		inputString    string
		expectedResult *Tree
	}{
		{
			name:           "Single parser matches",
			parsers:        Sequence(digitParser),
			inputString:    "123",
			expectedResult: &Tree{MatchedRunes: []rune("123"), Subtree: []*Tree{
				{MatchedRunes: []rune("123"), Tag: "digits"},
			}},
		},
		{
			name:           "Both parsers match",
			parsers:        Sequence(digitParser, letterParser).WithTag("test"),
			inputString:    "123abc",
			expectedResult: &Tree{MatchedRunes: []rune("123abc"), Tag: "test", Subtree: []*Tree{
				{MatchedRunes: []rune("123"), Tag: "digits", MatchStartIndex: 0},
				{MatchedRunes: []rune("abc"), Tag: "letters", MatchStartIndex: 3},
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
			expectedResult: &Tree{MatchedRunes: []rune("")},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seqParser := tc.parsers
			inputRunes := []rune(tc.inputString)
			result := seqParser.Parse(inputRunes, 0, NewCache())
			test.Eq(t, tc.expectedResult, result)
		})
	}
}