package util

import (
	"sync"
)

type StringTreeSet struct {
	mux   *sync.Mutex
	nodes map[string]*StringTreeSet
	value interface{}
}

//NewStringTreeSet creates a new datastructure based on hash tree
func NewStringTreeSet() *StringTreeSet {
	tree := new(StringTreeSet)
	var mux sync.Mutex
	tree.nodes = make(map[string]*StringTreeSet)
	tree.mux = &mux
	return tree
}

func newStringTreeSetNode(parent *StringTreeSet, value interface{}) *StringTreeSet {
	tree := new(StringTreeSet)
	tree.nodes = make(map[string]*StringTreeSet)
	tree.mux = parent.mux
	tree.value = value
	return tree
}

//Add node to tree
func (tree *StringTreeSet) Add(key string, value interface{}) *StringTreeSet {
	tree.mux.Lock()
	defer tree.mux.Unlock()
	return tree.add(key, value)
}

func (tree *StringTreeSet) add(key string, value interface{}) *StringTreeSet {
	if tree.nodes == nil {
		tree.nodes = make(map[string]*StringTreeSet)
	}
	tree.nodes[key] = newStringTreeSetNode(tree, value)
	return tree.nodes[key]
}

func (tree *StringTreeSet) AddPath(value interface{}, path ...string) {
	tree.mux.Lock()
	defer tree.mux.Unlock()
	root := tree
	for _, p := range path {
		root = root.add(p, nil)
	}
	root.value = value
}

//Find a node by key path
func (tree *StringTreeSet) Find(keys ...string) *StringTreeSet {
	if len(keys) == 0 {
		return tree
	}
	node, ok := tree.nodes[keys[0]]
	if !ok {
		return nil
	}
	return node.Find(keys[1:]...)
}

//Exist path in a tree
func (tree *StringTreeSet) Exist(keys ...string) bool {
	return tree.Find(keys...) != nil
}

//Value returns node value
func (tree *StringTreeSet) Value() interface{} {
	return tree.value
}
