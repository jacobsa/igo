// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The deps package offers utility functions for dealing with dependency
// graphs.
package deps

import (
	"container/vector"
	"igo/set"
)

// BuildTotalOrder accepts a map from package names to the dependencies of
// those packages, and returns a safe order in which to compile them (assuming
// there are no circular dependencies). The result will contain only those
// packages which were present as keys in deps.
func BuildTotalOrder(deps map[string]*set.StringSet) []string {
	var visitor topologicalSortVisitor
	visitor.nodes = make(map[string]*packageNode)
	visitor.edges = deps

	for key, _ := range deps {
		visitor.nodes[key] = &packageNode{}
	}

	for key, _ := range deps {
		visitor.Visit(key)
	}

	return visitor.result.Data()
}

type packageNode struct {
	visited bool
}

// Implements a depth-first search topological sort algorithm for directed
// acyclic graphs.
type topologicalSortVisitor struct {
	result vector.StringVector
	nodes  map[string]*packageNode
	edges  map[string]*set.StringSet
}

func (v *topologicalSortVisitor) Visit(name string) {
	// Is this an unvisited node for a package we care about?
	node, ok := v.nodes[name]
	if !ok || node.visited {
		return
	}

	node.visited = true
	for otherName := range v.edges[name].Iter() {
		v.Visit(otherName)
	}

	v.result.Push(name)
}
