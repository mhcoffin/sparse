package speg

import "github.com/google/uuid"

type StarParser struct {
	id     uuid.UUID
	parser Parser
}

func (z StarParser) Omit() Parser {
	return NewOmitParser(z)
}

func (z StarParser) Star() StarParser {
	return z
}

func (z StarParser) Flatten() TokenParser {
	return newTokenParser(z)
}

func (z StarParser) Parse(input []rune, start int, ctx *Context) *Tree {
	pos := start
	var children []*Tree
	_, isOmitParser := z.parser.(OmitParser)
	for {
		child := z.parser.Parse(input, pos, ctx)
		if child == nil || len(child.Match) == 0 || pos == len(input) {
			result := &Tree{
				Start:    start,
				Match:    input[start:pos],
				Children: children,
			}
			ctx.setCachedValue(z.id, start, result)
			return result
		} else {
			pos += len(child.Match)
			if ctx.withChildren && !isOmitParser {
				children = append(children, child)
			}
		}
	}
}

func (z StarParser) ID() uuid.UUID {
	return z.id
}

func (z StarParser) Tagged(tag string) TaggedParser {
	return Tagged(z, tag)
}

// Star creates a parser that matches zero or more of p. 
func Star(p Parser) Parser {
	switch pp := p.(type) {
	case StarParser:
		return pp
	case Matcher:
		return pp.Star()
	default:
		return StarParser{
			id:     uuid.New(),
			parser: pp,
		}
	}
}
