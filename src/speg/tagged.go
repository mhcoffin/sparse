package speg

import "github.com/google/uuid"

type TaggedParser struct {
	id     uuid.UUID
	parser Parser
	tag    string
}

func (t TaggedParser) Star() Parser {
	return newStarParser(t)
}

func (t TaggedParser) Flatten() Parser {
	return newFlatParser(t)
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

func (t TaggedParser) Tag() string {
	return t.tag
}

func (t TaggedParser) Tagged(tag string) Parser {
	return TaggedParser{
		id:     uuid.New(),
		parser: t.parser,
		tag:    tag,
	}
}

// newTaggedParser returns a new parser that parses with parser and then attaches tag to the result.
func newTaggedParser(parser Parser, tag string) TaggedParser {
	return TaggedParser{
		id:     uuid.New(),
		parser: parser,
		tag:    tag,
	}
}
