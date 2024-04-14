package speg

import "github.com/google/uuid"

type SequenceParser struct {
	id         uuid.UUID
	subParsers []Parser
}

func (p SequenceParser) Star() Parser {
	return newStarParser(p)
}

func (p SequenceParser) Flatten() Parser {
	return newFlatParser(p)
}

func (p SequenceParser) ID() uuid.UUID {
	return p.id
}

func (p SequenceParser) Id() uuid.UUID {
	return p.id
}

func (p SequenceParser) Tag() string {
	return ""
}

func (p SequenceParser) Tagged(tag string) Parser {
	return TaggedParser{
		id:     uuid.New(),
		parser: p,
		tag:    tag,
	}
}

// Parse matches a sequence of parsers.
func (p SequenceParser) Parse(input []rune, start int, context *Context) *Tree {
	myResult, haveResult := context.getCachedValue(p.id, start)
	if haveResult {
		return myResult
	}

	if err := context.pushActive(p, start); err != nil {
		return nil
	}
	defer context.popActive()

	var position = start
	var children []*Tree

	for _, parser := range p.subParsers {
		result := parser.Parse(input, position, context)
		if result == nil {
			return nil
		}
		position += len(result.Match)
		if context.withChildren {
			children = append(children, result)
		}
	}
	result := &Tree{
		Match:    input[start:position],
		Start:    start,
		Children: children,
	}
	context.setCachedValue(p.id, start, result)
	return result
}

// Seq returns a parser that succeeds if all of its sub-parsers succeed in sequence.
func Seq(parsers ...Parser) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: parsers,
	}
}
