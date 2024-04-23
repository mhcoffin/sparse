package speg

import "github.com/google/uuid"

type NotParser struct {
	id ID
	parser Parser
}

func (n NotParser) Omit() Parser {
	return Omit(n)
}

func (n NotParser) Parse(input []rune, start int, ctx *Context) *Tree {
	x := n.parser.Parse(input, start, ctx)
	if x == nil {
		return &Tree{
			Start:    start,
		}
	} else {
		return nil
	}
}

func (n NotParser) ID() uuid.UUID {
	return n.id
}

// Not matches the empty string if p fails to match. If p matches
// Not(p) fails.
func Not(p Parser) NotParser {
	return NotParser{
		id: uuid.New(),
		parser: p,
	}
}

