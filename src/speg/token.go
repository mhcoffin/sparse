package speg

import (
	"github.com/google/uuid"
	"unicode"
)

type TokenParser struct {
	id     ID
	parser Parser
	tag    string
}

func (f TokenParser) Omit() Parser {
	return Omit(f)
}

func (f TokenParser) ID() uuid.UUID {
	return f.id
}

func (f TokenParser) Parse(input []rune, start int, ctx *Context) *Tree {
	pos := start
	for ; pos < len(input); pos++ {
		if !unicode.IsSpace(input[pos]) {
			break
		}
	}
	t := f.parser.Parse(input, pos, ctx.WithoutChildren())
	if t == nil {
		return nil
	}
	return &Tree{
		Start: start,
		Match: input[start:pos +len(t.Match)],
		Children: []*Tree{t},
		Tag: f.tag,
	}
}

func (f TokenParser) Star() Parser {
	return Star(f)
}

func (f TokenParser) Tagged(tag string) TaggedParser {
	return Tagged(f, tag)
}

// Token returns a parser that matches optional leading whitespace followed
// by whatever parser matches. If it succeeds, it will have a single child
// that contains the runes that parser matched. Moreover, the child will not
// itself have any children.  
func Token(parser Parser) TokenParser {
	return TokenParser{
		id:     uuid.New(),
		parser: parser,
	}
}
