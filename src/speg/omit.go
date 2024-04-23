package speg

import "github.com/google/uuid"

type OmitParser struct {
	id     ID
	tag string
	parser Parser
}

func (o OmitParser) Parse(input []rune, start int, ctx *Context) *Tree {
	result := o.parser.Parse(input, start, ctx.WithoutChildren())
	if result == nil {
		return nil
	}
	result.Omit = true
	return result
}

func (o OmitParser) ID() uuid.UUID {
	return o.id
}

func (o OmitParser) Omit() Parser {
	return o
}

func (o OmitParser) Tagged(tag string) OmitParser {
	return OmitParser{
		id: uuid.New(),
		tag: tag,
		parser: o.parser,
	}
}

func Omit(parser Parser) OmitParser {
	return OmitParser{
		id: uuid.New(),
		parser: parser,
	}
}
