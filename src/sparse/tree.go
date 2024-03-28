package sparse

// A Tree describes the result of applying a parser to an input.
//
// [Tree.Runes] contains the prefix of the input was matched by the parser.
// It may be empty.
//
// [Tree.Children] contains tagged subtrees. E.g.,
//
//	Seq(Digits.ParserBase("before"), Exactly("."), Digits.ParserBase("after"))
//
// returns will return a tree with two children: the first for digits tagged "before",
// the second for digits tagged "after".
type Tree struct {
	Runes    []rune
	Children []*Tree
	Tag      string
}

// String returns all the runes matched. It returns "" if the tree is nil.
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return string(t.Runes)
}

// Get returns the first node (in depth-first, left-to-right order) that has the
// specified tag, or nil if there is none. If t is nil, it returns nil.
func (t *Tree) Get(tag string) *Tree {
	if t == nil {
		return nil
	}
	if t.Tag == tag {
		return t
	}
	for _, child := range t.Children {
		if result := child.Get(tag); result != nil {
			return result
		}
	}
	return nil
}

// GetAll returns all nodes that have the specified tag, or nil if there are none.
func (t *Tree) GetAll(tag string) []*Tree {
	if t == nil {
		return nil
	}
	var result []*Tree
	if t.Tag == tag {
		result = append(result, t)
	}
	for _, child := range t.Children {
		result = append(result, child.GetAll(tag)...)
	}
	return result
}

// GetChild returns the first child (left-to-right order) that has the specified
// tag, or nil if there is none. If t is nil, it returns nil.
func (t *Tree) GetChild(tag string) *Tree {
	if t == nil {
		return nil
	}
	for _, child := range t.Children {
		if child.Tag == tag {
			return child
		}
	}
	return nil
}
