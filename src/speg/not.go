package speg

import "github.com/google/uuid"

type NotParser struct {
	id ID
	parser Parser
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

func Not(p Parser) NotParser {
	return NotParser{
		id: uuid.New(),
		parser: p,
	}
}

