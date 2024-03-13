package sparse

import (
	"strings"
	"unicode"
)

// Package sparse provides types and functions to implement simple parsers.
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
// A combinator is a parser that takes other parsers as parameters. E.g., [Optional]
// takes a parser as parameter and returns a new parser that tries to match using
// the parameter, but if that fails, succeeds with the empty prefix. Another example
// is [Seq]: Seq(Digits, Exactly("."), Digits) matches a sequence of one or more digits
// followed by a decimal point followed by more digits. [Seq] returns a tree
// that has one child for each parameter. (But see [Parser.Elide].)

// Parser is a function that matches a prefix of input and returns a tree that indicates
// what was matched.
type Parser func(input []rune) *Tree

// A Tree describes the result of applying a parser to an input. 
// 
// [Tree.Runes] is set to whatever prefix of the input was matched by the parser. It may be empty.
//
// [Tree.Children] is set to the results of sub-parsers. E.g., 
//   Seq(Digits, Exactly("."), Digits)
// returns a tree with three children: the first for the digits, the second for ".", and
// the third for more digits.
//
// The Tag is supplied by the user to make processing the parse tree easier. Any parser
// can be [Tagged]: Digits.Tagged(17) or Seq(Digits, Exactly("."), Digits).Tagged(9).
type Tree struct {
	Runes    []rune
	Children []*Tree
	Tag      int
}

const (
	// ElideTag is a special tag that indicates that a Tree should be elided in a sequence.
	ElideTag = -1
	FlattenTag = -2
)

// String returns all the runes matched in creating t, including elisions.
func (t *Tree) String() string {
	return string(t.Runes)
}

// Any matches any single rune. Not(Any) is a good way to match the end of input.
var Any Parser = func(input []rune) *Tree {
	if len(input) == 0 {
		return nil
	}
	return &Tree{Runes: input[:1]}
}

// Tagged returns a new parser that matches exactly what m matches. However, when m
// matches, the parse tree is tagged with the specified tag instead of the Default.
// Negative tags are reserved.
func (p Parser) Tagged(tag int) Parser {
	return func(input []rune) *Tree {
		tree := p(input)
		if tree == nil {
			return tree
		}
		tree.Tag = tag
		return tree
	}
}

// LookingAt creates a new parser that matches the empty prefix if p matches input.
// [LookingAt] is always Elided.  In regexp parlance, LookingAt is positive lookahead:
//   Seq(X, LookingAt(Y))
// is a parser that matches X, but only if it's followed by Y. 
func LookingAt(p Parser) Parser {
	return func(input []rune) *Tree {
		t := p(input)
		if t == nil {
			return nil
		}
		return &Tree{Tag: ElideTag, Runes: input[:0]}
	}
}

// Not returns a new parser that fails if p matches, and matches an empty
// slice if m fails. [Not] is automatically elided. In regexp parlance, Not 
// implements negative lookahead:
//   Seq(X, Not(Y))
// is a parser that matches X, but only if it's not followed by Y.
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

// Elide returns a new parser that matches p and marks the result as "elided".
// Some combinators, such as [Seq], [OneOrMore], [ZeroOrMore] omit elided results
// from their Children.  
// 
// E.g., if
//	Seq(Exactly("("), Digits, Exactly(")")) 
// succeeds, it returns a tree with three children. On the other hand,
//	Seq(Exactly("(").Elide(), Digits, Exactly(")").Elide()) 
// returns a tree with only one child, omitting the parentheses.
func (p Parser) Elide() Parser {
	return p.Tagged(ElideTag)
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

// WS matches zero or more whitespace runes. The result is always elided.
var WS Parser = func(input []rune) *Tree {
	pos := 0
	for pos < len(input) && unicode.IsSpace(input[pos]) {
		pos++
	}
	return &Tree{Tag: ElideTag, Runes: input[:pos]}
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

// NonEmpty returns a parser that matches m if m matches a non-empty prefix, otherwise fails.
// NonEmpty is pointless unless m can match an empty prefix.
func (p Parser) NonEmpty() Parser {
	return func(input []rune) *Tree {
		t := p(input)
		if t == nil || len(t.Runes) == 0 {
			return nil
		}
		return &Tree{Runes: t.Runes}
	}
}

// Seq returns a Parser that matches a sequence of parsers, left-to-right.
// If it succeeds, the result will have one child for each parser that matches
// and is not elided (see also [Parser.Elide]).
//
// If parsers is empty, Seq() succeeds with an empty prefix. 
func Seq(parsers ...Parser) Parser {
	return func(input []rune) *Tree {
		pos := 0
		children := make([]*Tree, 0, len(parsers))
		for _, m := range parsers {
			t := m(input[pos:])
			if t == nil {
				return nil
			}
			pos += len(t.Runes)
			if t.Tag != ElideTag {
				children = append(children, t)
			}
		}
		return &Tree{Runes: input[:pos], Children: children}
	}
}

// FirstOf returns a parser that tries a series of parsers one after another until one
// succeeds. The result is the result of the first parser that succeeds.
// If the parsers all fail, the resulting parser fails.
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
// The result will have one child for every match. E.g., 
//   ZeroOrMore(FirstOf(Exactly("foo"), Exactly("bar"))) ("foofoobar")
// will have three children matching ["foo", "foo", "bar"]
//
// Zero or more applied to a series of parsers matches the entire series zero or
// more times. E.g.,
//   ZeroOrMore(Exactly("("), Digits, Exactly(")")) ("(123)(0)(123x)")
// will match with six children: ["(", "123", ")", "(", "0", ")"]. The ending "(123x)"
// does not match the entire series of parsers, so parsing stops before that.
//
// Any of the parameters that are elided [Parser.Elide] will be omitted from the result. 
// E.g.,
//   ZeroOrMore(Exactly("(").Elide(), Digits, Exactly(")").Elide()) ("(123)(0)(123x)")
// will omit the parentheses and match with two children: ["123", "0"].
//   
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
				if tree.Tag != ElideTag {
					seq = append(seq, tree)
				}
			}
			pos = seqPos
			children = append(children, seq...)
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
					if len(children) == 0 {
						return nil
					}
					return &Tree{Runes: input[:pos], Children: children}
				}
				seqPos += len(tree.Runes)
				if tree.Tag != ElideTag {
					seq = append(seq, tree)
				}
			}
			pos = seqPos
			children = append(children, seq...)
		}
	}

}
