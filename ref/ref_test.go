package ref

import (
	"fmt"
	"io/ioutil"
	"testing"

	"golang.org/x/exp/rand"
)

func TestSet(t *testing.T) {
	g := NewGroup()
	for i := 0; i < 3; i++ {
		g.AddRoot(g.Add())
	}
	for i := 0; i < 30; i++ {
		g.Add()
	}

	fmt.Println(g.All())
	fmt.Println(g.Root())

	g.All().Each(func(from Id) bool {
		g.All().RandSubSet(rand.Int() % 3).Each(func(to Id) bool {
			g.Link(from, to)
			return true
		})
		return true
	})
	fmt.Println(g.Dot())

	ioutil.WriteFile("t.dot", []byte(g.Dot()), 0600)
}
