package speg

import (
	"github.com/google/uuid"
)

type ID = uuid.UUID

// A Parser is a thing that matches its input beginning at start and returns a Tree
// describing the match, or nil if there is no match.
type Parser interface {
	// Parse matches input[start:] and returns a Tree describing the result, or
	// nil if there is no match.
	Parse(input []rune, start int, ctx *Context) *Tree

	// ID returns the unique ID of this parser.
	ID() uuid.UUID

	// Tag returns the tag of this parser.
	Tag() string

	// Tagged returns a modified parser that tags trees with the specified tag.
	Tagged(tag string) Parser

	// Flatten returns a modified parser that doesn't allocate children.
	Flatten() Parser
	
	// Star returns a modified parser that matches zero or more times.
	Star() Parser
}

