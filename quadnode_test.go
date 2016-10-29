package rquad

import (
	"testing"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

func testNodeUniqueID(t *testing.T, fn newQuadtreeFunc) {
	var (
		img *binimg.Binary
		err error
	)

	img, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)

	scanner, err := imgscan.NewScanner(img)
	check(t, err)
	q, err := fn(scanner, 1)
	check(t, err)

	ids := make(map[int]struct{})
	q.ForEachLeaf(Gray, func(n Node) {
		if _, ok := ids[n.ID()]; ok {
			t.Errorf("id %d already exists", n.ID())
		}
		ids[n.ID()] = struct{}{}
	})
}

func TestBasicNodeUniqueID(t *testing.T) {
	testNodeUniqueID(t, newBasicTree)
}

func TestCNNodeUniqueID(t *testing.T) {
	testNodeUniqueID(t, newCNTree)
}
