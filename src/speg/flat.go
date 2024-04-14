package speg

import "github.com/google/uuid"

// A FlatParser parses without allocating any children.
type FlatParser struct {
	id     ID
	parser Parser
	tag    string
}

func (f FlatParser) Star() Parser {
	return newStarParser(f)
}

func (f FlatParser) Parse(input []rune, start int, ctx *Context) *Tree {
	return f.Parse(input, start, ctx.WithoutChildren())
}

func (f FlatParser) ID() uuid.UUID {
	return f.id
}

func (f FlatParser) Tag() string {
	return f.tag
}

func (f FlatParser) Tagged(tag string) Parser {
	return FlatParser{id: f.id, parser: f.parser, tag: tag}
}

func (f FlatParser) Flatten() Parser {
	return f.parser
}

func newFlatParser(parser Parser) FlatParser {
	return FlatParser{
		id:     uuid.New(),
		parser: parser,
	}
}
