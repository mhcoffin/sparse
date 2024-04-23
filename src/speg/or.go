package speg

import "github.com/google/uuid"

type OrParser struct {
	id         uuid.UUID
	subParsers []Parser
}

func (p OrParser) Omit() Parser {
	return Omit(p)
}

// Star is a convenience: p.Star() is equivalent to Star(p).
func (p OrParser) Star() Parser {
	return Star(p)
}

func (p OrParser) ID() uuid.UUID {
	return p.id
}

// Tagged is a convenience method equivalent to Tagged(p, tag)
func (p OrParser) Tagged(tag string) TaggedParser {
	return Tagged(p, tag)
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

// Or returns a new parser that tries each of parsers in turn and returns results
// from the first one that succeeds. If none succeed it fails and return nil.
func Or(parsers ...Parser) Parser {
	return OrParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}
