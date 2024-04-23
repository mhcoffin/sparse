package speg

import "github.com/google/uuid"

type IndirectParser struct {
	parser **Parser
}

func (d IndirectParser) Parse(input []rune, start int, ctx *Context) *Tree {
	if d.parser == nil {
		panic("Indirect parser used before definition")
	}
	return (**d.parser).Parse(input, start, ctx)
}

func (d IndirectParser) Omit() Parser {
	return Omit(d)
}

// ID returns the ID of the underlying parser. If the underlying
// parser is not set, returns uuid.Nil (all zeros)
func (d IndirectParser) ID() uuid.UUID {
	if *d.parser == nil {
		panic("Indirect parser used before definition")
	}
	return (**d.parser).ID()
}

// Indirect creates a proxy parser that delegates all operations to *p. This enables
// recursive parsers. E.g., 
// 
//   var expr Parser
//   expr = Seq(Exactly("("), Indirect(&expr), Exactly(")"))
// 
// It is an error to invoke ID() or Parse() on an indirect parser before p is defined.
func Indirect(p *Parser) IndirectParser {
	return IndirectParser{
		parser: &p,
	}
}
