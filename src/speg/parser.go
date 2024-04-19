package speg

import (
	"github.com/google/uuid"
)

type ID = uuid.UUID

// A Parser matches its input beginning at start and returns a Tree
// describing the match, or nil if there is no match.
type Parser interface {
	// Parse matches input[start:] and returns a Tree describing the result, or
	// nil if there is no match.
	Parse(input []rune, start int, ctx *Context) *Tree
	
	// ID returns the unique ID of this parser.
	ID() uuid.UUID
	
	// Omit returns a new parser whose result is omitted from Children in the result. 
	Omit() Parser
}

