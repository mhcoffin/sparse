package speg

import (
	"github.com/google/uuid"
	"unicode"
)

// A Tree describes the result of matching some input.
type Tree struct {
	Start    int
	Match    []rune
	Children []*Tree
	Tag      string
}

type ID = uuid.UUID

// A Cache holds Trees previously produced for this input.
// It is indexed first by the ID of the parser, then by the input location.
type Cache = map[uuid.UUID]map[int]*Tree

func NewCache() Cache {
	return make(map[ID]map[int]*Tree)
}

// A Parser is a thing that matches its input beginning at start and returns a Tree
// describing the match, or nil if there is no match.
type Parser interface {
	// Parse matches input[start:] and returns a Tree describing the result, or
	// nil if it fails to match.
	Parse(input []rune, start int, cache Cache) *Tree

	// ID returns the unique ID for this parser.
	ID() uuid.UUID

	// Tag returns the tag of this parser.
	Tag() string

	// Tagged returns a new parser (with a new ID) that tags its results with tag.
	Tagged(tag string) Parser
}

// A MatchingFunc is a function that tries to match a prefix of its input.
// If it succeeds, it returns a Tree with Runes set to the prefix that
// it matches. If it fails, it returns nil.
type MatchingFunc = func(input []rune) *Tree

type TagField struct {
	tag string
}

// A Matcher is a parser that employs a MatchingFunc to scan input.
type Matcher struct {
	TagField
	id           ID
	matchingFunc MatchingFunc
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
	tree := m.matchingFunc(input[start:])
	if tree == nil {
		myCache[start] = nil
		return nil
	}
	tree.Start = start
	tree.Tag = m.tag
	myCache[start] = tree
	return tree
}

func (m Matcher) ID() uuid.UUID {
	return m.id
}

func (m Matcher) Tag() string {
	return m.tag
}

// Tagged returns a new Matcher that parses exactly what m does, but
// tags result trees with tag.
func (m Matcher) Tagged(tag string) Parser {
	return Matcher{
		TagField:     TagField{tag: tag},
		id:           m.id,
		matchingFunc: m.matchingFunc,
	}
}

// NewMatcher creates a matchingFunc from a MatcherFunc.
// Distinct calls to NewMatcher produce different IDs.
func NewMatcher(m MatchingFunc) Matcher {
	return Matcher{
		id:           uuid.New(),
		matchingFunc: m,
	}
}

// Any matches any single rune. It fails only if the input is empty.
func Any() Matcher {
	result := NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 {
			return nil
		}
		return &Tree{
			Match: input[:1],
		}
	})
	return result
}

// Letter matches any single unicode letter. If fails if the first
// rune in the input is not a letter.
func Letter() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 || !unicode.IsLetter(input[0]) {
			return nil
		}
		return &Tree{
			Match: input[:1],
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
					Match: input[:k],
				}
			}
		}
		if len(input) == 0 {
			return nil
		}
		return &Tree{Match: input}
	})
}

// Digit matches any single unicode digit.
func Digit() Matcher {
	return NewMatcher(func(input []rune) *Tree {
		if len(input) == 0 || !unicode.IsDigit(input[0]) {
			return nil
		}
		return &Tree{
			Match: input[:1],
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
					Match: input[:k],
				}
			}
		}
		if len(input) == 0 {
			return nil
		}
		return &Tree{Match: input}
	})
}

type FirstOfParser struct {
	id         uuid.UUID
	subParsers []Parser
	TagField
}

func (p FirstOfParser) ID() uuid.UUID {
	return p.id
}

func (p FirstOfParser) Tag() string {
	return p.tag
}

func (p FirstOfParser) Tagged(tag string) Parser {
	return FirstOfParser{
		id:         uuid.New(),
		subParsers: p.subParsers,
		TagField:   TagField{tag: tag},
	}
}

func (p FirstOfParser) Parse(input []rune, start int, cache Cache) *Tree {
	myCache := getCache(cache, p.id)
	myResult, haveResult := myCache[start]
	if haveResult {
		return myResult
	}
	for _, parser := range p.subParsers {
		try := parser.Parse(input, start, cache)
		if try != nil {
			if try.Tag == "" {
				result := &Tree{
					Start: start,
					Match: try.Match,
					Tag:   p.tag,
				}
				myCache[start] = result
				return result
			} else {
				result := &Tree{
					Start:    start,
					Match:    try.Match,
					Children: []*Tree{try},
					Tag:      p.tag,
				}
				myCache[start] = result
				return result
			}
		}
	}
	myCache[start] = nil
	return nil
}

// FirstOf returns the result from the first of parsers that succeeds, or
// nil if none of them do. If it succeeds, it will have one child: the match
// that succeeded.
func FirstOf(parsers ...Parser) FirstOfParser {
	return FirstOfParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}

type SequenceParser struct {
	id         uuid.UUID
	subParsers []Parser
	TagField
}

func (p SequenceParser) ID() uuid.UUID {
	return p.id
}

// Id returns the ID of this parser.
func (p SequenceParser) Id() uuid.UUID {
	return p.id
}

func (p SequenceParser) Tag() string {
	return p.tag
}

// Tagged returns a new SequenceParser that tags its results with tag.
func (p SequenceParser) Tagged(tag string) Parser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: p.subParsers,
		TagField:   TagField{tag},
	}
}

// Parse matches a sequence of parsers.
func (p SequenceParser) Parse(input []rune, start int, cache Cache) *Tree {
	var position = start
	var subtrees []*Tree
	for _, parser := range p.subParsers {
		result := parser.Parse(input, position, cache)
		if result == nil {
			return nil
		}
		position += len(result.Match)
		if result.Tag != "" {
			subtrees = append(subtrees, result)
		}
	}

	return &Tree{
		Match:    input[start:position],
		Start:    start,
		Children: subtrees,
		Tag:      p.tag,
	}
}

// Sequence returns a parser that succeeds if all of its sub-parsers succeed in sequence.
func Sequence(parsers ...Parser) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}

type OptionalParser struct {
	id     uuid.UUID
	parser Parser
	TagField
}

func (o OptionalParser) Parse(input []rune, start int, cache Cache) *Tree {
	tree := o.parser.Parse(input, start, cache)
	if tree == nil {
		return &Tree{
			Start:    start,
			Match:    input[:0],
			Children: nil,
			Tag:      o.tag,
		}
	} else if tree.Tag != "" {
		return &Tree{
			Start:    start,
			Match:    tree.Match,
			Children: []*Tree{tree},
			Tag:      o.tag,
		}
	} else {
		return &Tree{
			Start: start,
			Match: tree.Match,
			Tag:   o.tag,
		}
	}
}

func (o OptionalParser) Tagged(tag string) Parser {
	return &OptionalParser{
		id:       uuid.New(),
		parser:   o.parser,
		TagField: TagField{tag: tag},
	}
}

func (o OptionalParser) ID() uuid.UUID {
	return o.id
}

func (o OptionalParser) Tag() string {
	return o.tag
}

// Opt optionally matches parser.
func Opt(parser Parser) Parser {
	return OptionalParser{
		id:     uuid.New(),
		parser: parser,
	}
}
