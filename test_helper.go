package rquad

import (
	"testing"

	"github.com/arl/imgtools/imgscan"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Fatal(err)
	}
}

type newQuadtreeFunc func(imgscan.Scanner, int) (Quadtree, error)

func newBasicTree(scanner imgscan.Scanner, resolution int) (Quadtree, error) {
	return NewBasicTree(scanner, resolution)
}

func newCNTree(scanner imgscan.Scanner, resolution int) (Quadtree, error) {
	return NewCNTree(scanner, resolution)
}

func neighbourColors(n Node) (white, black int) {
	ForEachNeighbour(n, func(nb Node) {
		switch nb.Color() {
		case Black:
			black++
		case White:
			white++
		}
	})
	return
}
