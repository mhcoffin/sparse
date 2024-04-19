package speg

import "github.com/google/uuid"

type TokenParser struct {
	id     ID
	parser Parser
	tag    string
}

func (f TokenParser) Omit() Parser {
	return f
}

func (f TokenParser) ID() uuid.UUID {
	return f.id
}

func (f TokenParser) Parse(input []rune, start int, ctx *Context) *Tree {
	return f.Parse(input, start, ctx.WithoutChildren())
}

func (f TokenParser) Star() Parser {
	return Star(f)
}

func (f TokenParser) Tagged(tag string) TaggedParser {
	return Tagged(f, tag)
}

// newTokenParser creates a parser that matches exactly what parser matches,
// but avoids allocating children. The result of using a TokenParser is a Tree
// without any substructure.
func newTokenParser(parser Parser) TokenParser {
	return TokenParser{
		id:     uuid.New(),
		parser: parser,
	}
}
