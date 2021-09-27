// Package doublejump provides a revamped Google's jump consistent hash.
package doublejump

import (
	"math/rand"

	"github.com/fufuok/balancer/internal/go-jump"
)

type looseHolder struct {
	a []interface{}
	m map[interface{}]int
	f []int
}

func (l *looseHolder) add(obj interface{}) {
	if _, ok := l.m[obj]; ok {
		return
	}

	if nf := len(l.f); nf == 0 {
		l.a = append(l.a, obj)
		l.m[obj] = len(l.a) - 1
	} else {
		idx := l.f[nf-1]
		l.f = l.f[:nf-1]
		l.a[idx] = obj
		l.m[obj] = idx
	}
}

func (l *looseHolder) remove(obj interface{}) {
	if idx, ok := l.m[obj]; ok {
		l.f = append(l.f, idx)
		l.a[idx] = nil
		delete(l.m, obj)
	}
}

func (l *looseHolder) get(key uint64) interface{} {
	na := len(l.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key, na)
	return l.a[h]
}

func (l *looseHolder) shrink() {
	if len(l.f) == 0 {
		return
	}

	var a []interface{}
	for _, obj := range l.a {
		if obj != nil {
			a = append(a, obj)
			l.m[obj] = len(a) - 1
		}
	}
	l.a = a
	l.f = nil
}

type compactHolder struct {
	a []interface{}
	m map[interface{}]int
}

func (c *compactHolder) add(obj interface{}) {
	if _, ok := c.m[obj]; ok {
		return
	}

	c.a = append(c.a, obj)
	c.m[obj] = len(c.a) - 1
}

func (c *compactHolder) shrink(a []interface{}) {
	for i, obj := range a {
		c.a[i] = obj
		c.m[obj] = i
	}
}

func (c *compactHolder) remove(obj interface{}) {
	if idx, ok := c.m[obj]; ok {
		n := len(c.a)
		c.a[idx] = c.a[n-1]
		c.m[c.a[idx]] = idx
		c.a[n-1] = nil
		c.a = c.a[:n-1]
		delete(c.m, obj)
	}
}

func (c *compactHolder) get(key uint64) interface{} {
	na := len(c.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key*0xc6a4a7935bd1e995, na)
	return c.a[h]
}

// Hash is a revamped Google's jump consistent hash. It overcomes the shortcoming of the
// original implementation - not being able to remove nodes.
type Hash struct {
	loose   looseHolder
	compact compactHolder
}

// NewHash creates a new doublejump hash instance, which does NOT threadsafe.
func NewHash() *Hash {
	hash := &Hash{}
	hash.loose.m = make(map[interface{}]int)
	hash.compact.m = make(map[interface{}]int)
	return hash
}

// Add adds an object to the hash.
func (h *Hash) Add(obj interface{}) {
	if obj == nil {
		return
	}

	h.loose.add(obj)
	h.compact.add(obj)
}

// Remove removes an object from the hash.
func (h *Hash) Remove(obj interface{}) {
	if obj == nil {
		return
	}

	h.loose.remove(obj)
	h.compact.remove(obj)
}

// Len returns the number of objects in the hash.
func (h *Hash) Len() int {
	return len(h.compact.a)
}

// LooseLen returns the size of the inner loose object holder.
func (h *Hash) LooseLen() int {
	return len(h.loose.a)
}

// Shrink removes all empty slots from the hash.
func (h *Hash) Shrink() {
	h.loose.shrink()
	h.compact.shrink(h.loose.a)
}

// Get returns an object according to the key provided.
func (h *Hash) Get(key uint64) interface{} {
	obj := h.loose.get(key)
	switch obj {
	case nil:
		return h.compact.get(key)
	default:
		return obj
	}
}

// All returns all the objects in this Hash.
func (h *Hash) All() []interface{} {
	n := len(h.compact.a)
	if n == 0 {
		return nil
	}
	all := make([]interface{}, n)
	copy(all, h.compact.a)
	return all
}

// Random returns a random object.
func (h *Hash) Random() interface{} {
	if n := len(h.compact.a); n > 0 {
		idx := rand.Intn(n)
		return h.compact.a[idx]
	}
	return nil
}
