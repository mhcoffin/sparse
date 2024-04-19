package speg

import "github.com/google/uuid"

type SequenceParser struct {
	id         uuid.UUID
	subParsers []Parser
}

func (p SequenceParser) Omit() Parser {
	return NewOmitParser(p)
}

func (p SequenceParser) Star() Parser {
	return Star(p)
}

func (p SequenceParser) Flatten() TokenParser {
	return newTokenParser(p)
}

func (p SequenceParser) ID() uuid.UUID {
	return p.id
}

func (p SequenceParser) Tagged(tag string) TaggedParser {
	return Tagged(p, tag)
}



// Parse matches a sequence of parsers, left to right. The result Tree will have one
// child for each of the parsers.
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
		
		_, isOmitParser := parser.(OmitParser)
		if context.withChildren && !isOmitParser {
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

// Seq returns a parser that succeeds if all of its sub-parsers succeed, left-to-right.
// The resulting Tree will have one child for each non-omitted sub-parsers.
func Seq(subParsers ...Parser) SequenceParser {
	return SequenceParser{
		id:         uuid.New(),
		subParsers: subParsers,
	}
}
