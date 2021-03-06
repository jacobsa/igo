// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The set package offers an unordered, unique container for strings, like
// mathematical sets. Currently it is not at all optimized for speed.
//
// TODO(jacobsa): This probably belongs in the built-in container/ package set.
package set

import (
	"container/vector"
)

// A StringSet is a container of strings whose fundamental operations are
// insertion and testing for membership.
type StringSet struct {
	elements vector.StringVector
}

// Contains returns true or false based on whether the supplied string is in
// the set.
func (set *StringSet) Contains(s string) bool {
	for val := range set.elements.Iter() {
		if val == s {
			return true
		}
	}

	return false
}

// Insert adds s to the set, so that Contains(s) will now be true if it wasn't
// already.
func (set *StringSet) Insert(s string) {
	if !set.Contains(s) {
		set.elements.Push(s)
	}
}

// Union adds the elements from s to d.
func (d *StringSet) Union(s *StringSet) {
	for val := range s.Iter() {
		d.Insert(val)
	}
}

// Iter returns an iterator for the elements in the set. The order of the
// elements is not guaranteed.
func (set *StringSet) Iter() <-chan string { return set.elements.Iter() }
