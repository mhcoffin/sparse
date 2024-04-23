package speg

import "github.com/google/uuid"

type LeftRecursiveParser struct {
	id           uuid.UUID
	base         Parser
	continuation Parser
	tag          string
}

func (l LeftRecursiveParser) Omit() Parser {
	return Omit(l)
}

func (l LeftRecursiveParser) Star() Parser {
	return Star(l)
}

func (l LeftRecursiveParser) Flatten() TokenParser {
	return Token(l)
}

func (l LeftRecursiveParser) Parse(input []rune, start int, ctx *Context) *Tree {
	base := l.base.Parse(input, start, ctx)
	if base == nil {
		return nil
	}
	pos := start + len(base.Match)

	cont := l.continuation.Parse(input, pos, ctx)
	if cont == nil || len(cont.Match) == 0 {
		return base
	}
	var children []*Tree
	if ctx.withChildren {
		children = append(children, base)
		children = append(children, cont.Children...)
	}
	pos += len(cont.Match)
	lhs := &Tree{
		Start:    start,
		Match:    input[start:pos],
		Children: children,
		Tag:      l.tag,
	}
	for {
		cont := l.continuation.Parse(input, pos, ctx)
		if cont == nil || len(cont.Match) == 0 {
			// TODO: add tag?
			return lhs
		}
		pos += len(cont.Match)
		var children []*Tree
		if ctx.withChildren {
			children = append(children, lhs)
			children = append(children, cont.Children...)
		}
		lhs = &Tree{
			Start:    start,
			Match:    input[start:pos],
			Children: children,
			Tag:      l.tag,
		}
	}
}

func (l LeftRecursiveParser) ID() uuid.UUID {
	return l.id
}

func (l LeftRecursiveParser) Tagged(tag string) Parser {
	return LeftRecursiveParser{
		id:           uuid.New(),
		base:         l.base,
		continuation: l.continuation,
		tag:          tag,
	}
}

func Left(base Parser, continuation Parser) LeftRecursiveParser {
	return LeftRecursiveParser{
		id:           uuid.New(),
		base:         base,
		continuation: continuation,
	}
}
