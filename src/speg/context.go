package speg

import (
	"fmt"
	"github.com/google/uuid"
)

type Context struct {
	cache         Cache
	activeParsers []ActiveParser
	withChildren  bool
}

// A Cache holds Trees previously produced for this input.
// It is indexed first by the ID of the parser, then by the input location.
type Cache = map[uuid.UUID]map[int]*Tree

type ActiveParser struct {
	id    ID
	start int
}

func (context *Context) pushActive(parser Parser, start int) error {
	if context.isActive(parser, start) {
		return fmt.Errorf("left recursion")
	}
	context.activeParsers = append(context.activeParsers, ActiveParser{parser.ID(), start})
	return nil
}

func (context *Context) popActive() {
	context.activeParsers = context.activeParsers[:len(context.activeParsers)-1]
}

func (context *Context) isActive(parser Parser, start int) bool {
	for k := len(context.activeParsers) - 1; k >= 0; k-- {
		if context.activeParsers[k].start == start && context.activeParsers[k].id == parser.ID() {
			return true
		}
	}
	return false
}

func (context *Context) getCache(id ID) map[int]*Tree {
	parserCache, ok := context.cache[id]
	if !ok {
		parserCache = make(map[int]*Tree)
		context.cache[id] = parserCache
	}
	return parserCache
}

func (context *Context) getCachedValue(id ID, pos int) (*Tree, bool) {
	cache := context.getCache(id)
	value, ok := cache[pos]
	return value, ok
}

func (context *Context) setCachedValue(id ID, pos int, value *Tree) {
	context.getCache(id)[pos] = value
}

func (context *Context) WithoutChildren() *Context {
	return &Context{
		cache:         context.cache,
		activeParsers: context.activeParsers,
		withChildren:  false,
	}
}

func NewContext() *Context {
	return &Context{
		cache:         make(map[ID]map[int]*Tree),
		activeParsers: []ActiveParser{},
		withChildren:  true,
	}
}
