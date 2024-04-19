package speg

import "github.com/google/uuid"

type OmitParser struct {
	id     ID
	parser Parser
}

func (o OmitParser) Parse(input []rune, start int, ctx *Context) *Tree {
	result := o.parser.Parse(input, start, ctx.WithoutChildren())
	result.Omit = true
	return result
}

func (o OmitParser) ID() uuid.UUID {
	return o.id
}

func (o OmitParser) Omit() Parser {
	return o
}

func NewOmitParser(parser Parser) OmitParser {
	return OmitParser{
		id: uuid.New(),
		parser: parser,
	}
}
