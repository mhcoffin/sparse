package speg

import "github.com/google/uuid"

type IndirectParser struct {
	parser **Parser
}

func (d *IndirectParser) Parse(input []rune, start int, ctx *Context) *Tree {
	if d.parser == nil {
		panic("Indirect.Parse() with nil parser")
	}
	result := (**d.parser).Parse(input, start, ctx)
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
	return ""
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

// Indirect creates a parser that invokes another parser, to be set later.
// It is used in recursive definitions:
//
//	iTerm := Indirect()
//	term := Seq(iTerm, Exactly("*"), iTerm)
//	iTerm.Set(&term)
//
// Now term is a parser that recursively parses two sub-terms separated by "*".
func Indirect() *IndirectParser {
	var ptr *Parser
	return &IndirectParser{
		parser: &ptr,
	}
}
