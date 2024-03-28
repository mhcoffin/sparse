package speg

import (
	"github.com/google/uuid"
	"unicode"
)

// A Tree describes the result of matching some input.
type Tree struct {
	MatchedRunes    []rune
	MatchStartIndex int
	Subtree         []*Tree
	Tag             string
}

type ID = uuid.UUID

// A Cache holds Trees previously produced for this input.
// It is indexed first by the ID of the parser, then by the input location.
type Cache = map[uuid.UUID]map[int]*Tree

func NewCache() Cache {
	return make(map[ID]map[int]*Tree )
}

// A Parser is a thing that matches its input beginning at start and returns a Tree
// describing the match, or nil if there is no match.
type Parser interface {
	// Parse matches input[start:] and returns a Tree describing the result, or
	// nil if it fails to match.
	Parse(input []rune, start int, cache Cache) *Tree

	// GetMatchID returns the match ID for this parser. Two parsers with the same
	// Match ID must produce the same result.
	GetMatchID() uuid.UUID
}

// A MatchingFunc is a function that tries to match a prefix of its input.
// If it succeeds, it returns a Tree with Runes set to the prefix that
// it matches. If it fails, it returns nil.
type MatchingFunc = func(input []rune) *Tree

type TagField struct {
	tag string
}

// A Matcher is a parser that employs a MatchingFunc to scan input and tags
// its result Tree with tag.
type Matcher struct {
	TagField
	id      ID
	matcher MatchingFunc
}

// getCache returns the cache for Matcher, creating a new one if necessary.
func getCache(cache Cache, id ID) map[int]*Tree {
	myCache, ok := cache[id]
	if !ok {
		myCache = make(map[int]*Tree)
		cache[id] = myCache
	}
	return myCache
}

func (m Matcher) Parse(input []rune, start int, cache Cache) *Tree {
	myCache := getCache(cache, m.id)
	cachedResult, isCached := myCache[start]
	if isCached {
		return cachedResult
	}
	tree := m.matcher(input[start:])
	if tree == nil {
		myCache[start] = nil
		return nil
	}
	tree.MatchStartIndex = start
	tree.Tag = m.tag
	myCache[start] = tree
	return tree
}

func (m Matcher) GetMatchID() uuid.UUID {
	return m.id
}

// WithTag returns a new matcher that tags result trees with tag.
func (m Matcher) WithTag(tag string) Matcher {
	return Matcher{
		TagField: TagField{tag: tag},
		id:       m.id,
		matcher:  m.matcher,
	}
}

// NewMatcher creates a matcher from a MatcherFunc.
// Distinct calls to NewMatcher produce different IDs.
func NewMatcher(m MatchingFunc) Matcher {
	return Matcher{
		id:      uuid.New(),
		matcher: m,
	}
}

// Any matches any single rune. It fails only if the input is empty.
func Any() Matcher {
	result := NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 {
			return nil
		}
		return &Tree{
			MatchedRunes: input[:1],
		}
	})
	return result
}

func Letter() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 || !unicode.IsLetter(input[0]) {
			return nil
		}
		return &Tree{
			MatchedRunes: input[:1],
		}
	})
}

// OneOrMoreLetters matches one or more unicode letter runes.
func OneOrMoreLetters() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		for k, r := range input {
			if !unicode.IsLetter(r) {
				if k == 0 {
					return nil
				}
				return &Tree{
					MatchedRunes: input[:k],
				}
			}
		}
		if len(input) == 0 {
			return nil
		}
		return &Tree{MatchedRunes: input}
	})
}

// Digit matches any single unicode digit.
func Digit() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 || !unicode.IsDigit(input[0]) {
			return nil
		}
		return &Tree{
			MatchedRunes: input[:1],
		}
	})
}

func OneOrMoreDigits() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		for k, r := range input {
			if !unicode.IsDigit(r) {
				if k == 0 {
					return nil
				}
				return &Tree{
					MatchedRunes: input[:k],
				}
			}
		}
		if len(input) == 0 {
			return nil
		}
		return &Tree{MatchedRunes: input}
	})
}

type FirstOfParser struct {
	id         uuid.UUID
	subParsers []Parser
}

func (p FirstOfParser) Id() uuid.UUID {
	return p.id
}

func (p FirstOfParser) Parse(input []rune, start int, cache Cache) *Tree {
	myCache := getCache(cache, p.id)
	myResult, haveResult := myCache[start]
	if haveResult {
		return myResult
	}
	for _, parser := range p.subParsers {
		result := parser.Parse(input, start, cache)
		if result != nil {
			myCache[start] = result
			return result
		}
	}
	myCache[start] = nil
	return nil
}

// FirstOf returns the result from the first of parsers that succeeds, or
// nil if none of them do. FirstOf does not create its own parser tree. It
// just directly returns the result from parsers.
func FirstOf(parsers ...Parser) FirstOfParser {
	return FirstOfParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}

type SequenceParser struct {
	id         uuid.UUID
	subParsers []Parser
	tag        string
}

func (p SequenceParser) GetMatchID() uuid.UUID {
	return p.id
}

// Id returns the ID of this parser.
func (p SequenceParser) Id() uuid.UUID {
	return p.id
}

// Tagged returns a new SequenceParser that tags its results with tag.
func (p SequenceParser) Tagged(tag string) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: p.subParsers,
		tag:        tag,
	}
}

// Parse matches the input starting from start, using the sub-parsers of the SequenceParser.
// It returns a Tree describing the result if the match is successful, or nil otherwise.
// The result Tree contains the matched runes, the starting index of the match,
// and any sub-parsers that tag their results.
// The top-level tree will be tagged with the p's tag.
func (p SequenceParser) Parse(input []rune, start int, cache Cache) *Tree {
	var position = start
	var subtrees []*Tree
	for _, parser := range p.subParsers {
		result := parser.Parse(input, position, cache)
		if result == nil {
			return nil
		}
		position += len(result.MatchedRunes)
		if result.Tag != "" {
			subtrees = append(subtrees, result)
		}
	}

	return &Tree{
		MatchedRunes:    input[start:position],
		MatchStartIndex: start,
		Subtree:         subtrees,
		Tag: p.tag,
	}
}

// Sequence returns a parser that succeeds if all of its subparsers succeed in sequence.
func Sequence(parsers ...Parser) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}

func (p SequenceParser) WithTag(tag string) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: p.subParsers,
		tag:        tag,
	}
}

