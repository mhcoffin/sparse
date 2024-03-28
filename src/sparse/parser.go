package sparse

import (
	"strings"
	"unicode"
)

// Package sparse provides types and functions to implement simple parser.
//
// A Parser is a function that takes a slice of runes and either succeeds or
// fails to match a prefix (possibly empty) of the slice. If it succeeds it
// returns a tree that describes what it matched  If it fails it returns nil.
//
// A primitive parser matches a prefix as a single unit. E.g., [Any] matches any
// single rune, and [Digits] matches a sequence of one or more digits. A primitive
// parser returns a tree that indicates the portion of the input it matched,
// but without any children.
//
// A combinator is a parser that takes other parser as parameters. E.g., [Optional]
// takes a parser as parameter and returns a new parser that tries to match using
// the parameter, but if that fails, succeeds with the empty prefix. Another example
// is [Seq]: Seq(Digits, Exactly("."), Digits) matches a sequence of one or more digits
// followed by a decimal point followed by more digits. [Seq] returns a tree
// that has one child for each parameter. (But see [Parser.Elide].)

// A Parser is a function that either fails and returns nil,
// or matches a prefix of input and returns a tree that indicates
// what was matched.
type Parser func(input []rune) *Tree

// Any returns a parser that matches any single rune.
var Any Parser = func(input []rune) *Tree {
	if len(input) == 0 {
		return nil
	}
	return &Tree{Runes: input[:1]}
}

// OneOf matches any of the runes in s.
func OneOf(s string) Parser {
	// TODO: maybe use a bit map if the runes are all small.
	m := make(map[rune]bool)
	for _, c := range s {
		m[c] = true
	}
	return func(input []rune) *Tree {
		if len(input) == 0 {
			return nil
		}
		_, ok := m[input[0]]
		if ok {
			return &Tree{Runes: input[:1]}
		}
		return nil
	}
}

// ZeroOrMoreOf returns a parser that matches the longest prefix that consists
// of runes that are in s.
func ZeroOrMoreOf(s string) Parser {
	// TODO: maybe use a bit map if the runes are all small.
	m := make(map[rune]bool)
	for _, c := range s {
		m[c] = true
	}
	return func(input []rune) *Tree {
		pos := 0
		for pos < len(input) {
			if _, ok := m[input[pos]]; !ok {
				return &Tree{Runes: input[:pos]}
			}
			pos++
		}
		return &Tree{Runes: input}
	}
}

// Tagged returns a new parser that matches exactly what p matches. However, if p
// succeeds, the resulting parse tree will be tagged with the specified tag instead
// of the default ("").
func (p Parser) Tagged(tag string) Parser {
	if p == nil {
		panic("nil parser in ParserBase")
	}

	return func(input []rune) *Tree {
		tree := p(input)
		if tree == nil {
			return tree
		}
		tree.Tag = tag
		return tree
	}
}

// LookingAt returns a parser that matches the empty prefix if p matches the input.
// LookingAt is used for positive lookahead:
//
//	Seq(X, LookingAt(Y))
//
// is a parser that matches X, but only if it's followed by Y.
func LookingAt(p Parser) Parser {
	return func(input []rune) *Tree {
		t := p(input)
		if t == nil {
			return nil
		}
		return &Tree{Runes: input[:0]}
	}
}

// Not returns a new parser that fails if p matches, and succeeds for an
// empty prefix if m fails. Not is mostly used for negative lookahead:
//
//	Seq(X, Not(Y))
//
// is a parser that matches X, but only if it's not followed by Y.
//
//	Not(Any)
//
// matches the end of input.
func Not(p Parser) Parser {
	return func(input []rune) *Tree {
		tree := p(input)
		if tree == nil {
			return &Tree{Runes: input[:0]}
		} else {
			return nil
		}
	}
}

// Digit is a Parser that matches any single Unicode digit.
var Digit Parser = func(input []rune) *Tree {
	if len(input) > 0 && unicode.IsDigit(input[0]) {
		return &Tree{Runes: input[:1]}
	}
	return nil
}

// Digits matches one or more unicode digits.
var Digits Parser = func(input []rune) *Tree {
	if len(input) == 0 || !unicode.IsDigit(input[0]) {
		return nil
	}
	pos := 1
	for pos < len(input) && unicode.IsDigit(input[pos]) {
		pos++
	}
	return &Tree{Runes: input[:pos]}
}

// Letter is a Parser that matches any single Unicode letter.
var Letter Parser = func(input []rune) *Tree {
	if len(input) > 0 && unicode.IsLetter(input[0]) {
		return &Tree{Runes: input[:1]}
	}
	return nil
}

// Letters is a Parser that matches one more unicode letters.
var Letters Parser = func(input []rune) *Tree {
	if len(input) == 0 || !unicode.IsLetter(input[0]) {
		return nil
	}
	pos := 1
	for pos < len(input) && unicode.IsLetter(input[pos]) {
		pos++
	}
	return &Tree{Runes: input[:pos]}
}

// Space is a Parser that matches any single Unicode space character.
var Space Parser = func(input []rune) *Tree {
	if len(input) > 0 && unicode.IsSpace(input[0]) {
		return &Tree{Runes: input[:1]}
	}
	return nil
}

// Spaces matches one or more unicode whitespace characters.
var Spaces Parser = func(input []rune) *Tree {
	if len(input) == 0 || !unicode.IsSpace(input[0]) {
		return nil
	}
	pos := 1
	for pos < len(input) && unicode.IsSpace(input[pos]) {
		pos++
	}
	return &Tree{Runes: input[:pos]}
}

// WS matches zero or more whitespace runes.
var WS Parser = func(input []rune) *Tree {
	pos := 0
	for pos < len(input) && unicode.IsSpace(input[pos]) {
		pos++
	}
	return &Tree{Runes: input[:pos]}
}

// Exactly returns a parser that succeeds when s is a prefix of the input.
func Exactly(s string) Parser {
	return func(input []rune) *Tree {
		var length = 0
		for index, r := range s {
			if index >= len(input) || input[index] != r {
				return nil
			}
			length++
		}
		return &Tree{Runes: input[:length]}
	}
}

// IgnoreCase returns a parser that succeeds when s is a prefix of input ignoring case.
// [unicode.ToLower] is used to canonicalize case.
func IgnoreCase(s string) Parser {
	target := strings.ToLower(s)
	return func(input []rune) *Tree {
		var length = 0
		for index, r := range target {
			if index >= len(input) || unicode.ToLower(input[index]) != r {
				return nil
			}
			length++
		}
		return &Tree{Runes: input[:length]}
	}
}

// Optional returns a parser that either matches m and returns that result, or matches the
// empty prefix and returns that.
func Optional(m Parser) Parser {
	return func(input []rune) *Tree {
		t := m(input)
		if t == nil {
			return &Tree{Runes: input[:0]}
		}
		return t
	}
}

// NonEmpty returns a parser that matches m if m matches a non-empty prefix,
// otherwise fails. NonEmpty is pointless unless m can match an empty prefix.
func (p Parser) NonEmpty() Parser {
	return func(input []rune) *Tree {
		t := p(input)
		if t == nil || len(t.Runes) == 0 {
			return nil
		}
		return &Tree{Runes: t.Runes}
	}
}

// Seq returns a [Parser] that matches a sequence of parsers, left-to-right.
// If it succeeds, the result will have one child for each of the parameters
// that are tagged. E.g., if
//
//	Seq(Digits.ParserBase("digits"), WS, Letters.ParserBase("letters")
//
// succeeds on some input, the result will have two children, the first tagged
// "digits" and the second "letters".
//
// If parser is empty, Seq() succeeds with an empty prefix.
func Seq(parsers ...Parser) Parser {
	for _, parser := range parsers {
		if parser == nil {
			panic("nil parser in Seq")
		}
	}
	return func(input []rune) *Tree {
		pos := 0
		children := make([]*Tree, 0, len(parsers))
		for _, m := range parsers {
			t := m(input[pos:])
			if t == nil {
				return nil
			}
			pos += len(t.Runes)
			if t.Tag != "" {
				children = append(children, t)
			}
		}
		return &Tree{Runes: input[:pos], Children: children}
	}
}

// FirstOf returns a parser that tries a series of parser one after another until one
// succeeds. The result is the result of the first parser that succeeds.
// If the parser all fail, the resulting parser fails.
func FirstOf(parsers ...Parser) Parser {
	return func(input []rune) *Tree {
		for _, m := range parsers {
			t := m(input)
			if t != nil {
				return t
			}
		}
		return nil
	}
}

// ZeroOrMore applied to a single parser matches that parser zero or more times.
// The result will have one child for every match that is [ParserBase]. E.g.,
//
//	ZeroOrMore(FirstOf(Letters.ParserBase("word"), Digits.ParserBase("number")))
//
// when applied to "first123second456" will have four children:
//
//	     {Tag: "word", Runes: []rune("first")},
//			{Tag: "number", Runes: []rune("123")},
//			{Tag: "word", Runes: []rune("second")},
//			{Tag: "number", Runes: []rune("456")},
//
// When applied to ",first123", it will match the empty prefix and have no children.
//
// Zero or more applied to a series of parsers matches the entire series zero or
// more times. E.g.,
//
//	ZeroOrMore(Exactly("number"), Digits.ParserBase("number"))
//
// when applied to "number17number11number" will have two children:
//
//	{Tag: "number", Runes: []rune("17")},
//	{Tag: "number", Runes: []rune("11")},
//
// The trailing "number" is left unmatched since it is not followed by digits.
func ZeroOrMore(parsers ...Parser) Parser {
	return func(input []rune) *Tree {
		pos := 0
		var children []*Tree
		for {
			var seq []*Tree
			seqPos := pos
			for _, parser := range parsers {
				tree := parser(input[seqPos:])
				if tree == nil {
					return &Tree{Runes: input[:pos], Children: children}
				}
				seqPos += len(tree.Runes)
				if tree.Tag != "" {
					seq = append(seq, tree)
				}
			}
			pos = seqPos
			if seq != nil {
				children = append(children, seq...)
			}
		}
	}
}

// OneOrMore is like ZeroOrMore but fails unless there is at least one match.
func OneOrMore(m ...Parser) Parser {
	return func(input []rune) *Tree {
		pos := 0
		var children []*Tree
		for {
			var seq []*Tree
			seqPos := pos
			for _, parser := range m {
				tree := parser(input[seqPos:])
				if tree == nil {
					if pos == 0 {
						return nil
					}
					return &Tree{Runes: input[:pos], Children: children}
				}
				seqPos += len(tree.Runes)
				if tree.Tag != "" {
					seq = append(seq, tree)
				}
			}
			pos = seqPos
			if seq != nil {
				children = append(children, seq...)
			}
		}
	}
}

// Deref provides a way to break dependency loops.
// E.g.,
//
//	var Foo Parser
//	Expr = FirstOf(
//	  Seq(Digits,
//	  Not(Exactly("+"))),
//	  Seq(Digits, Exactly("+"), Deref(&Expr)))
func Deref(x *Parser) Parser {
	return func(input []rune) *Tree {
		if x == nil {
			panic("Deref used with nil parser")
		}
		return (*x)(input)
	}
}


