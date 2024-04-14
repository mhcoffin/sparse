package speg

import (
	"fmt"
	"strconv"
	"strings"
)

// A Tree describes the result of matching some input.
type Tree struct {
	Start    int
	Match    []rune
	Children []*Tree
	Tag      string
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
