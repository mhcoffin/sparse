package sparse

import (
	"github.com/shoenig/test"
	"testing"
)

var tree = &Tree{
	Runes: []rune("top"),
	Tag:   "top",
	Children: []*Tree{
		{
			Tag:      "child 0",
			Runes:    []rune("child 0"),
			Children: nil,
		},
		{
			Tag:   "child 1",
			Runes: []rune("another child"),
			Children: []*Tree{
				{
					Tag:      "child 1/0",
					Runes:    []rune("child 1.0"),
					Children: nil,
				},
				{
					Tag:   "child 1/1",
					Runes: []rune("child 1.1"),
					Children: []*Tree{
						{
							Tag:      "duplicate",
							Runes:    []rune("duplicate 1"),
							Children: nil,
						},
						{
							Tag:      "duplicate",
							Runes:    []rune("duplicate 2"),
							Children: nil,
						},
					},
				},
			},
		},
		{
			Tag:   "child 2",
			Runes: []rune("child two"),
			Children: []*Tree{
				{
					Tag:      "child 2/0",
					Runes:    []rune("child 2.0"),
					Children: nil,
				},
				{
					Tag:      "child 2/1",
					Runes:    []rune("child 3.2"),
					Children: nil,
				},
				{
					Tag:   "child 2/2",
					Runes: []rune("child 2.2"),
					Children: []*Tree{
						{
							Tag:      "child 2/2/0",
							Runes:    []rune("child 2.2.0"),
							Children: nil,
						},
						{
							Tag:      "child 2/2/1",
							Runes:    []rune("child 2.2.1"),
							Children: nil,
						},
						{
							Tag:      "duplicate",
							Runes:    []rune("duplicate 3"),
							Children: nil,
						},
					},
				},
			},
		},
	},
}

func TestTree_Get(t *testing.T) {
	testCases := []struct {
		tag      string
		expected *Tree
	}{
		{tag: "top", expected: tree},
		{tag: "child 0", expected: tree.Children[0]},
		{tag: "child 1", expected: tree.Children[1]},
		{tag: "child 2", expected: tree.Children[2]},
		{tag: "child 0/0", expected: nil},
		{tag: "child 1/0", expected: tree.Children[1].Children[0]},
		{tag: "child 1/1", expected: tree.Children[1].Children[1]},
		{tag: "duplicate", expected: tree.Children[1].Children[1].Children[0]},
	}

	for _, tt := range testCases {
		t.Run(tt.tag, func(t *testing.T) {
			test.Eq(t, tree.Get(tt.tag), tt.expected)
		})
	}
}

func TestTree_GetAll(t *testing.T) {
	testCases := []struct {
		tag      string
		expected []*Tree
	}{
		{tag: "top", expected: []*Tree{tree}},
		{tag: "child 0", expected: []*Tree{tree.Children[0]}},
		{tag: "child 1", expected: []*Tree{tree.Children[1]}},
		{tag: "child 2", expected: []*Tree{tree.Children[2]}},
		{tag: "child 0/0", expected: nil},
		{tag: "child 1/0", expected: []*Tree{tree.Children[1].Children[0]}},
		{tag: "child 1/1", expected: []*Tree{tree.Children[1].Children[1]}},
		{tag: "child 2/2/0", expected: []*Tree{tree.Children[2].Children[2].Children[0]}},
		{tag: "duplicate", expected: []*Tree{
			tree.Children[1].Children[1].Children[0],
			tree.Children[1].Children[1].Children[1],
			tree.Children[2].Children[2].Children[2],
		}},
	}

	for _, tt := range testCases {
		t.Run(tt.tag, func(t *testing.T) {
			test.Eq(t, tree.GetAll(tt.tag), tt.expected)
		})
	}

}

func TestTree_GetChild(t *testing.T) {
	testCases := []struct {
		tag      string
		expected *Tree
	}{
		{tag: "top", expected: nil},
		{tag: "child 0", expected: tree.Children[0]},
		{tag: "child 1", expected: tree.Children[1]},
		{tag: "child 2", expected: tree.Children[2]},
		{tag: "child 3", expected: nil},
		{tag: "child 0/1", expected: nil},
	}

	for _, tt := range testCases {
		t.Run(tt.tag, func(t *testing.T) {
			test.Eq(t, tree.GetChild(tt.tag), tt.expected)
		})
	}
}

func TestTree_Nil(t *testing.T) {
	var tree *Tree = nil
	test.Eq(t, tree.Get("foo"), nil)
	test.Eq(t, tree.GetAll("foo"), nil)
	test.Eq(t, tree.GetChild("foo"), nil)
	test.Eq(t, tree.String(), "")
}

func TestTree_String(t *testing.T) {
	tests := []struct {
		tree *Tree
		expected string
	}{
		{tree: &Tree{Tag: "top", Runes: []rune("foobar")}, expected: "foobar"},
		{tree: &Tree{Tag: "top", Runes: []rune{}}, expected: ""},
		{tree: &Tree{Tag: "top"}, expected: ""},
	}
	for _, tt := range tests {
		t.Run(string(tt.tree.Runes), func(t1 *testing.T) {
			test.Eq(t, tt.expected, tt.tree.String())
		})
	}
}

func TestNonEmpty(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			
		})
	}
}

