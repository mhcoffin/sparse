package speg

import "github.com/google/uuid"

type OrParser struct {
	id         uuid.UUID
	subParsers []Parser
}

func (p OrParser) Star() Parser {
	return newStarParser(p)
}

func (p OrParser) Flatten() Parser {
	return newFlatParser(p)
}

func (p OrParser) ID() uuid.UUID {
	return p.id
}

func (p OrParser) Tag() string {
	return ""
}

func (p OrParser) Tagged(tag string) Parser {
	return TaggedParser{
		id:     uuid.New(),
		parser: p,
		tag:    tag,
	}
}

func (p OrParser) Parse(input []rune, start int, context *Context) *Tree {
	result, ok := context.getCachedValue(p.id, start)
	if ok {
		return result
	}

	for _, parser := range p.subParsers {
		try := parser.Parse(input, start, context)
		if try != nil {
			return try
		}
	}
	context.setCachedValue(p.id, start, nil)
	return nil
}

// Or returns the result from the first of parsers that succeeds, or
// nil if none of them do. If it succeeds, it will have one child: the match
// that succeeded.
func Or(parsers ...Parser) Parser {
	return OrParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}
