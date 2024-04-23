package speg

import "github.com/google/uuid"

type TaggedParser struct {
	id     uuid.UUID
	parser Parser
	tag    string
}

func (t TaggedParser) Omit() Parser {
	return Omit(t)
}

func (t TaggedParser) Parse(input []rune, start int, ctx *Context) *Tree {
	tree := t.parser.Parse(input, start, ctx)
	if tree == nil {
		return nil
	}
	tree.Tag = t.tag
	return tree
}

func (t TaggedParser) ID() uuid.UUID {
	return t.id
}

func Tagged(p Parser, tag string) TaggedParser {
	return TaggedParser{
		id:     uuid.New(),
		parser: p,
		tag:    tag,
	}
}
