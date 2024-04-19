package speg

import "github.com/google/uuid"

type OptionalParser struct {
	id     uuid.UUID
	parser Parser
}

func (o OptionalParser) Omit() Parser {
	return NewOmitParser(o)
}

func (o OptionalParser) Star() Parser {
	return Star(o.parser)
}

func (o OptionalParser) Token() TokenParser {
	return TokenParser{
		id:     uuid.New(),
		parser: o,
	}
}

func (o OptionalParser) Parse(input []rune, start int, context *Context) *Tree {
	myResult, haveResult := context.getCachedValue(o.id, start)
	if haveResult {
		return myResult
	}

	tree := o.parser.Parse(input, start, context)
	if tree == nil {
		tree = &Tree{
			Start: start,
		}
	}
	context.setCachedValue(o.id, start, tree)
	return tree
}

func (o OptionalParser) ID() uuid.UUID {
	return o.id
}

func (o OptionalParser) Tag() string {
	return ""
}

func (o OptionalParser) Tagged(tag string) TaggedParser {
	return Tagged(o, tag)
}

// Opt returns the result from parser if parser succeeds, or an empty
// match if parser fails.
func Opt(parser Parser) Parser {
	return OptionalParser{
		id:     uuid.New(),
		parser: parser,
	}
}
