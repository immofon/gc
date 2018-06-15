package ref

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type Id int

type Set struct {
	m map[Id]bool
}

func (set Set) Exist(id Id) bool {
	_, exist := set.m[id]
	return exist
}

func (set Set) Add(id Id) {
	set.m[id] = true
}

func (set Set) Remove(id Id) {
	delete(set.m, id)
}

func (set Set) Len() int {
	return len(set.m)
}

func (set Set) Each(fn func(Id) bool) {
	for id, _ := range set.m {
		ok := fn(id)
		if !ok {
			return
		}
	}
}

func (set Set) Union(o Set) Set {
	n := NewSet()
	add := func(id Id) bool {
		n.Add(id)
		return true
	}
	set.Each(add)
	o.Each(add)
	return n
}

func (set Set) Sub(from Set) Set {
	i := NewSet().Union(set)
	from.Each(func(id Id) bool {
		i.Remove(id)
		return true
	})
	return i
}

func (set Set) RandSubSet(max int) Set {
	if max > set.Len() {
		max = set.Len()
	}

	sub := NewSet()
	rand_bool := func() bool {
		return rand.Int()%10 == 0 // 1/10
	}

	for sub.Len() < max {
		set.Each(func(id Id) bool {
			if rand_bool() {
				sub.Add(id)
			}
			if sub.Len() != max {
				return true
			} else {
				return false
			}
		})
	}

	return sub
}

type Ref struct {
	set Set
}

type Group struct {
	refs   map[Id]Ref
	root   Set
	nextId Id
}

func NewSet() Set {
	return Set{
		m: make(map[Id]bool),
	}
}

func NewRef() Ref {
	return Ref{
		set: NewSet(),
	}
}

func (r Ref) Refs() Set {
	return r.set
}

func NewGroup() *Group {
	return &Group{
		refs:   make(map[Id]Ref),
		root:   NewSet(),
		nextId: 1,
	}
}

func (g *Group) Add() Id {
	var id = g.nextId
	g.nextId += 1

	g.refs[id] = NewRef()
	return id
}

func (g *Group) Assert(id Id) {
	if _, ok := g.refs[id]; !ok {
		panic(fmt.Sprintf("ref id:%d don't exist\n", id))
	}
}

func (g *Group) Get(id Id) Ref {
	g.Assert(id)
	return g.refs[id]

}

func (g *Group) AddRoot(id Id) Ref {
	ref := g.Get(id)
	g.root.Add(id)
	return ref
}

func (g *Group) Link(from, to Id) {
	g.Assert(from)
	g.Assert(to)

	from_ref := g.Get(from)

	from_ref.set.Add(to)
}

func (g *Group) All() Set {
	set := NewSet()
	for id, _ := range g.refs {
		set.Add(id)
	}
	return set
}

func (g *Group) Root() Set {
	return g.root
}

func (g *Group) Dot() string {
	w := bytes.NewBuffer(nil)
	rand.Seed(uint64(time.Now().UnixNano()))
	node := func(id Id) {
		fmt.Fprintf(w, "\tnode%d[label=%d];\n", id, id)
	}
	root := func(id Id) {
		fmt.Fprintf(w, "\tnode%d[style=filled,color=gold];\n", id)
	}
	scaned := func(id Id) {
		fmt.Fprintf(w, "\tnode%d[style=filled,color=grey];\n", id)
	}
	marked := func(id Id) {
		fmt.Fprintf(w, "\tnode%d[style=filled,color=lightpink];\n", id)
	}
	link := func(from, to Id) {
		fmt.Fprintf(w, "\tnode%d->node%d;\n", from, to)
	}

	fmt.Fprintln(w, "digraph {")
	fmt.Fprintln(w, "\tnode[shape=circle];")

	// node
	g.All().Each(func(id Id) bool {
		node(id)
		return true
	})

	// root node
	g.Root().Each(func(id Id) bool {
		root(id)
		return true
	})

	// scaned node
	g.Scaned().Sub(g.Root()).Each(func(id Id) bool {
		scaned(id)
		return true
	})

	// marked node
	g.Marked().Each(func(id Id) bool {
		marked(id)
		return true
	})

	// eage
	g.All().Each(func(from Id) bool {
		g.Get(from).Refs().Each(func(to Id) bool {
			link(from, to)
			return true
		})
		return true
	})
	fmt.Fprintln(w, "}")
	return w.String()
}

func (g *Group) Scaned() Set {
	unscaned := NewSet()
	scaned := NewSet()

	unscaned = unscaned.Union(g.Root())
	for unscaned.Len() > 0 {
		unscaned.Each(func(id Id) bool {
			scaned.Add(id)
			unscaned.Remove(id)
			refs := g.Get(id).Refs()
			scaned.Each(func(i Id) bool {
				refs.Remove(i)
				return true
			})
			unscaned = unscaned.Union(refs)
			return true
		})
	}

	return scaned
}

func (g *Group) Marked() Set {
	return g.All().Sub(g.Scaned())

}
