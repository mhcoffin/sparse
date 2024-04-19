package speg

import (
	"github.com/google/uuid"
	"unicode"
)

// A MatchingFunc is a function that tries to match a prefix of its input.
// If it succeeds, it returns the length of the match. If it fails, it returns -1.
type MatchingFunc = func(input []rune) int

// A Matcher is a parser that employs a MatchingFunc to scan input.
type Matcher struct {
	id           ID
	matchingFunc MatchingFunc
	tag          string
}

func (m Matcher) Star() Matcher {
	return Matcher{
		id: uuid.New(),
		matchingFunc: func(input []rune) int {
			result := 0
			for {
				length := m.matchingFunc(input[result:])
				if length <= 0 {
					return result
				}
				result += length
			}
		},
	}
}

func (m Matcher) Parse(input []rune, start int, ctx *Context) *Tree {
	cachedResult, isCached := ctx.getCachedValue(m.id, start)
	if isCached {
		return cachedResult
	}

	length := m.matchingFunc(input[start:])
	if length == -1 {
		ctx.setCachedValue(m.id, start, nil)
		return nil
	}
	result := &Tree{
		Start: start,
		Match: input[start : start+length],
		Tag:   m.tag,
	}
	ctx.setCachedValue(m.id, start, result)
	return result
}

func (m Matcher) Tag() string {
	return m.tag
}

func (m Matcher) ID() uuid.UUID {
	return m.id
}

func (m Matcher) Tagged(tag string) Parser {
	return Matcher{
		id:           uuid.New(),
		matchingFunc: m.matchingFunc,
		tag:          tag,
	}
}

func (m Matcher) Omit() Parser {
	return NewOmitParser(m)
}

// NewMatcher creates a Matcher from a MatcherFunc.
func NewMatcher(m MatchingFunc) Matcher {
	return Matcher{
		id:           uuid.New(),
		matchingFunc: m,
	}
}

// Any matches any single rune. It fails only if the input is empty.
func Any() Matcher {
	result := NewMatcher(func(input []rune) int {
		if len(input) == 0 {
			return -1
		}
		return 1
	})
	return result
}

// Letter matches any single unicode letter. It fails if the first
// rune in the input is not a letter.
func Letter() Matcher {
	return NewMatcher(func(input []rune) int {
		if len(input) == 0 || !unicode.IsLetter(input[0]) {
			return -1
		}
		return 1
	})
}

// Letters matches one or more unicode letter runes.
func Letters() Matcher {
	return NewMatcher(func(input []rune) int {
		if len(input) == 0 || !unicode.IsLetter(input[0]) {
			return -1
		}
		for k, r := range input {
			if !unicode.IsLetter(r) {
				return k
			}
		}
		return len(input)
	})
}

// Digit matches any single unicode digit.
func Digit() Matcher {
	return NewMatcher(func(input []rune) int {
		if len(input) == 0 || !unicode.IsDigit(input[0]) {
			return -1
		}
		return 1
	})
}

func Digits() Matcher {
	return NewMatcher(func(input []rune) int {
		if len(input) == 0 || !unicode.IsDigit(input[0]) {
			return -1
		}
		for k, r := range input {
			if !unicode.IsDigit(r) {
				return k
			}
		}
		return len(input)
	})
}

func Exactly(s string) Matcher {
	return NewMatcher(func(input []rune) int {
		for pos, r := range s {
			if pos >= len(input) || input[pos] != r {
				return -1
			}
		}
		return len(s)
	})
}

func WhiteSpace() Matcher {
	return NewMatcher(func(input []rune) int {
		if len(input) == 0 || !unicode.IsSpace(input[0]) {
			return -1
		}
		for pos, r := range input {
			if !unicode.IsSpace(r) {
				return pos
			}
		}
		return len(input)
	})
}
