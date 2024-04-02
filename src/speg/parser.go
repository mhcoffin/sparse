package speg

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"unicode"
)

// A Tree describes the result of matching some input.
type Tree struct {
	Start    int
	Match    []rune
	Children []*Tree
	Tag      string
}

func (t *Tree) String() string {
	if t == nil {
		return "<nil>"
	}
	if t.Children == nil {
		return fmt.Sprintf("(%s %s)", t.Tag, string(t.Match))
	}
	var children []string
	for _, s := range t.Children {
		children = append(children, s.String())
	}
	return fmt.Sprintf("(%s %s)", t.Tag, strings.Join(children, " "))
}

func (t *Tree) Matched() string {
	return string(t.Match)
}

type ID = uuid.UUID

type ActiveParser struct {
	id    ID
	start int
}

func (context *Context) pushActive(id ID, start int) error {
	if context.isActive(id, start) {
		return fmt.Errorf("left recursion")
	}
	context.activeParsers = append(context.activeParsers, ActiveParser{id, start})
	return nil
}

func (context *Context) popActive() {
	context.activeParsers = context.activeParsers[:len(context.activeParsers)-1]
}

func (context *Context) isActive(id ID, start int) bool {
	for k := len(context.activeParsers) - 1; k >= 0; k-- {
		if context.activeParsers[k].start == start && context.activeParsers[k].id == id {
			return true
		}
	}
	return false
}

func (context *Context) getCache(id ID) map[int]*Tree {
	parserCache, ok := context.cache[id]
	if !ok {
		parserCache = make(map[int]*Tree)
		context.cache[id] = parserCache
	}
	return parserCache
}

type Context struct {
	cache         Cache
	activeParsers []ActiveParser
}

func NewContext() Context {
	return Context{
		cache:         make(map[ID]map[int]*Tree),
		activeParsers: []ActiveParser{},
	}
}

// A Cache holds Trees previously produced for this input.
// It is indexed first by the ID of the parser, then by the input location.
type Cache = map[uuid.UUID]map[int]*Tree

// A Parser is a thing that matches its input beginning at start and returns a Tree
// describing the match, or nil if there is no match.
type Parser interface {
	// Parse matches input[start:] and returns a Tree describing the result, or
	// nil if it fails to match.
	Parse(input []rune, start int, ctx Context) *Tree

	// ID returns the unique ID for this parser.
	ID() uuid.UUID

	// Tag returns the tag of this parser.
	Tag() string

	String() string
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
	name         string
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

func (m Matcher) Parse(input []rune, start int, context Context) *Tree {
	myCache := getCache(context.cache, m.id)
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

func (m Matcher) String() string {
	if m.tag == "" {
		return m.name
	} else {
		return fmt.Sprintf("%s:%s", m.name, m.tag)
	}
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
	return &Matcher{
		TagField:     TagField{tag: tag},
		id:           m.id,
		matchingFunc: m.matchingFunc,
		name: m.name,
	}
}

// NewMatcher creates a Matcher from a MatcherFunc.
// Distinct calls to NewMatcher produce different IDs.
func NewMatcher(m MatchingFunc, name string) Matcher {
	return Matcher{
		id:           uuid.New(),
		matchingFunc: m,
		name:         name,
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
	}, "any")
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
	}, "letter")
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
	}, "letter+")
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
	}, "digit")
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
	}, "digit+")
}

func Exactly(s string) Matcher {
	return NewMatcher(func(input []rune) *Tree {
		pos := 0
		for _, r := range s {
			if pos >= len(input) || input[pos] != r {
				return nil
			}
			pos++
		}
		return &Tree{Match: input[:pos]}
	}, fmt.Sprintf("`%s`", s))
}

type FirstOfParser struct {
	TagField
	id         uuid.UUID
	subParsers []Parser
}

func (p FirstOfParser) String() string {
	var sub []string
	for _, s := range p.subParsers {
		sub = append(sub, s.String())
	}
	var name string
	if p.tag == "" {
		name = "first"
	} else {
		name = fmt.Sprintf("first:%s", p.tag)
	}


	return fmt.Sprintf("(%s %s)", name, strings.Join(sub, " "))
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

func (p FirstOfParser) Parse(input []rune, start int, context Context) *Tree {
	cache := context.getCache(p.id)
	result, ok := cache[start]
	if ok {
		return result
	}

	for _, parser := range p.subParsers {
		var try *Tree
		if context.isActive(parser.ID(), start) {
			try = nil
		} else {
			try = parser.Parse(input, start, context)
		}
		if try != nil {
			if try.Tag == "" {
				result := &Tree{
					Start: start,
					Match: try.Match,
					Tag:   p.tag,
				}
				cache[start] = result
				return result
			} else {
				result := &Tree{
					Start:    start,
					Match:    try.Match,
					Children: []*Tree{try},
					Tag:      p.tag,
				}
				cache[start] = result
				return result
			}
		}
	}
	cache[start] = nil
	return nil
}

// FirstOf returns the result from the first of parsers that succeeds, or
// nil if none of them do. If it succeeds, it will have one child: the match
// that succeeded.
func FirstOf(parsers ...Parser) FirstOfParser {
	return FirstOfParser{
		id:         uuid.New(),
		subParsers: parsers,
		TagField:   TagField{"or"},
	}
}

type SequenceParser struct {
	id         uuid.UUID
	subParsers []Parser
	TagField
}

func (p SequenceParser) String() string {
	var sub []string
	for _, s := range p.subParsers {
		sub = append(sub, s.String())
	}
	var name string
	if p.tag == "" {
		name = "seq"
	} else {
		name = fmt.Sprintf("seq:%s", p.tag)
	}
	return fmt.Sprintf("(%s %s)", name, strings.Join(sub, " "))
}

func (p SequenceParser) ID() uuid.UUID {
	return p.id
}

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
func (p SequenceParser) Parse(input []rune, start int, context Context) *Tree {
	cache := context.getCache(p.id)
	myResult, haveResult := cache[start]
	if haveResult {
		return myResult
	}
	if err := context.pushActive(p.id, start); err != nil {
		panic("left recursion in Sequence")
	}
	defer context.popActive()

	var position = start
	var subtrees []*Tree
	for _, parser := range p.subParsers {
		result := parser.Parse(input, position, context)
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
		TagField:   TagField{"seq"},
	}
}

type OptionalParser struct {
	id     uuid.UUID
	parser Parser
	TagField
}

func (o OptionalParser) String() string {
	return fmt.Sprintf("opt(%s)", o.parser.String())
}

func (o OptionalParser) Parse(input []rune, start int, context Context) *Tree {
	cache := context.getCache(o.id)
	myResult, haveResult := cache[start]
	if haveResult {
		return myResult
	}

	tree := o.parser.Parse(input, start, context)
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
		id:       uuid.New(),
		parser:   parser,
		TagField: TagField{"opt"},
	}
}

type IndirectParser struct {
	parser **Parser
	tag    string
}

func (d *IndirectParser) Parse(input []rune, start int, ctx Context) *Tree {
	if d.parser == nil {
		panic("Indirect.Parse() with nil parser")
	}
	result := (**d.parser).Parse(input, start, ctx)
	result.Tag = d.tag
	return result
}

// ID returns the ID of the underlying parser. If the underlying
// parser is not set, returns uuid.Nil (all zeros)
func (d *IndirectParser) ID() uuid.UUID {
	if *d.parser == nil {
		return uuid.Nil
	}
	return (**d.parser).ID()
}

// Tag returns the tag of this parser.
func (d *IndirectParser) Tag() string {
	return d.tag
}

// Tagged creates a new IndirectParser that parses using the same underlying
// parser as d, but with the specified tag.
func (d *IndirectParser) Tagged(tag string) *IndirectParser {
	return &IndirectParser{
		parser: d.parser,
		tag:    tag,
	}
}

// Set assigns the parser p to d. I.e., when d is invoked, it will parse
// using p. Note, however that p's tag is overridden with d's tag, so that
// it is possible to have multiple IndirectParsers that use the same
// underlying parser but tag their results differently.
func (d *IndirectParser) Set(p *Parser) {
	if d.parser == nil {
		panic("nil in indirect parser!?")
	}
	*d.parser = p
}

func (d *IndirectParser) String() string {
	p := **d.parser
	return fmt.Sprintf("[&%s]", p.Tag())
}

// Indirect creates a parser that invokes another parser, to be set later.
// It is used in recursive definitions:
//
//	iTerm := Indirect()
//	term := Sequence(iTerm, Exactly("*"), iTerm)
//	iTerm.Set(&term)
//
// Now term is a parser that recursively parses two sub-terms separated by "*".
func Indirect() *IndirectParser {
	var ptr *Parser
	return &IndirectParser{
		parser: &ptr,
	}
}
