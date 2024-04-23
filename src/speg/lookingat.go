package speg

import "github.com/google/uuid"

type LookingAtParser struct {
	id     ID
	parser Parser
}

func (p LookingAtParser) Parse(input []rune, start int, ctx *Context) *Tree {
	x := p.parser.Parse(input, start, ctx.WithoutChildren())
	if x == nil {
		return nil
	}
	return &Tree{
		Start: start,
	}
}

func (p LookingAtParser) ID() uuid.UUID {
	return p.id
}

func (p LookingAtParser) Omit() Parser {
	return Omit(p)
}

// LookingAt returns a new parser that matches the empty string
// if parser would succeed, and fails if parser would fail. 
func LookingAt(parser Parser) LookingAtParser {
	return LookingAtParser{
		id:     uuid.New(),
		parser: parser,
	}
}
