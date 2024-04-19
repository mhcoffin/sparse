package speg

import (
	"fmt"
	"strconv"
	"strings"
)

// A Tree describes the result of matching input.
type Tree struct {
	// Position where the match starts
	Start    int
	// Slice of runes matched. May be empty. 
	Match    []rune
	// Children in match. May be empty (nil).
	Children []*Tree
	// User-specified tag of this tree.
	Tag      string
	// If Omit is true, this tree will be omitted from Children
	Omit bool
}

func (t *Tree) String() string {
	if t == nil {
		return "<nil>"
	}
	if t.Children == nil && t.Tag == "" {
		return fmt.Sprintf("%s", strconv.Quote(string(t.Match)))
	} else if t.Children == nil && t.Tag != "" {
		return fmt.Sprintf("(%s %s)", t.Tag, strconv.Quote(string(t.Match)))
	}

	var children []string
	for _, s := range t.Children {
		children = append(children, s.String())
	}
	if t.Tag == "" {
		return fmt.Sprintf("(%s)", strings.Join(children, " "))
	} else {
		return fmt.Sprintf("(%s %s)", t.Tag, strings.Join(children, " "))
	}
}

func (t *Tree) Matched() string {
	return string(t.Match)
}
