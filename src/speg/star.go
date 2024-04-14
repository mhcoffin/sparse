package speg

import "github.com/google/uuid"

type StarParser struct {
	id     uuid.UUID
	parser Parser
}

func (z StarParser) Star() Parser {
	return z
}

func (z StarParser) Flatten() Parser {
	return newFlatParser(z)
}

func (z StarParser) Parse(input []rune, start int, ctx *Context) *Tree {
	pos := start
	var children []*Tree
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
			if ctx.withChildren {
				children = append(children, child)
			}
		}
	}
}

func (z StarParser) ID() uuid.UUID {
	return z.id
}

func (z StarParser) Tag() string {
	return ""
}

func (z StarParser) Tagged(tag string) Parser {
	return TaggedParser{
		id:     uuid.New(),
		parser: z,
		tag:    tag,
	}
}

func newStarParser(p Parser) StarParser {
	return StarParser{
		id:     uuid.New(),
		parser: p,
	}
}
