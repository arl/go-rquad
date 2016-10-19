package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func TestBUQuadtreeLogicalErrors(t *testing.T) {
	var testTbl = []struct {
		w, h     int  // bitmap dimensions
		res      int  // resolution
		expected bool // true if the quadtree can be created without errors
	}{
		{0, 0, 2, false},   // 0x0 bitmap
		{10, 10, 0, false}, // resolution of 0
		{10, 10, 1, true},  // ok
		{10, 10, 6, false}, // can't be subdivided at least once
		{10, 13, 6, false}, // one dimension can't be subdivided at least once
		{13, 10, 6, false}, // one dimension can't be subdivided at least once
		{12, 12, 6, true},  // ok: 12 >= 6 * 2
		{13, 13, 6, true},  // ok: 13 >= 6 * 2
	}

	var (
		err     error
		bm      image.Image
		scanner binimg.Scanner
	)

	for _, tt := range testTbl {
		bm = binimg.New(image.Rect(0, 0, tt.w, tt.h))
		scanner, err = binimg.NewScanner(bm)
		check(t, err)
		_, err = NewBUQuadtree(scanner, tt.res)
		actual := err == nil
		if actual != tt.expected {
			t.Errorf("(%d,%d,%d): expected %v, actual %v, err:'%v'",
				tt.w, tt.h, tt.res, tt.expected, actual, err)
		}
	}
}

func TestBUQuadtreeSubdivisions(t *testing.T) {
	// this is a simple 32x32 image, white background with 3 black squares,
	// located so that they fill a quadrant the biggest of which is 8x8 pixels,
	// meaning that no nodes can ever be smaller than 8x8, that's why every
	// resolutions lower or equal than 8 should produce the same number of
	// nodes.
	var (
		err     error
		bm      image.Image
		scanner binimg.Scanner
	)
	bm, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	scanner, err = binimg.NewScanner(bm)
	check(t, err)

	for _, res := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		q, err := NewBUQuadtree(scanner, res)
		check(t, err)

		var whiteNodes QNodeList
		q.ForEachLeaf(White, appendNode(&whiteNodes))
		if len(whiteNodes) != 7 {
			t.Errorf("resolution:%d, expected 7 nodes, got %d",
				res, len(whiteNodes))
		}
	}

	for _, res := range []int{9, 15} {
		q, err := NewBUQuadtree(scanner, res)
		check(t, err)

		var whiteNodes QNodeList
		q.ForEachLeaf(White, appendNode(&whiteNodes))
		if len(whiteNodes) != 1 {
			t.Errorf("resolution:%d, expected 1 nodes, got %d",
				res, len(whiteNodes))
		}
	}
}

func TestBUQuadtreeQuery(t *testing.T) {
	var testTbl = []struct {
		pt     image.Point // queried point
		exists bool        // node exists
		eqRef  bool        // node should be equal to ref node
	}{
		{image.Pt(8, 0), true, true},
		{image.Pt(15, 0), true, true},
		{image.Pt(8, 7), true, true},
		{image.Pt(15, 7), true, true},
		{image.Pt(7, 0), true, false},
		{image.Pt(16, 0), true, false},
		{image.Pt(8, 8), true, false},
		{image.Pt(16, 8), true, false},
		{image.Pt(1, 31), true, false},
		{image.Pt(-1, 0), false, false},
		{image.Pt(32, 0), false, false},
		{image.Pt(0, -1), false, false},
		{image.Pt(0, 32), false, false},
	}

	// coordinate of reference node
	refPt := image.Pt(8, 0)
	var (
		err     error
		bm      image.Image
		scanner binimg.Scanner
	)
	bm, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	scanner, err = binimg.NewScanner(bm)
	check(t, err)

	for _, res := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		q, err := NewBUQuadtree(scanner, res)
		check(t, err)

		for _, tt := range testTbl {
			node := q.PointLocation(tt.pt)
			exists := node != nil
			if exists != tt.exists {
				t.Fatalf("resolution %d, expected exists to be %t for point %v, got %t instead",
					res, tt.exists, tt.pt, exists)
			}
			// obtain the refNode
			refNode := q.PointLocation(refPt)
			if refNode == nil {
				t.Fatalf("reference node should exist")
			}

			if exists {
				// check if queried node should be equal to reference node
				eqRef := node == refNode
				if eqRef != tt.eqRef {
					t.Errorf("resolution %d, got %t for comparison between reference node and queried node, expected %t. node:%v",
						res, eqRef, tt.eqRef, node)
				}
			}
		}
	}
}

func TestBUQuadtreeRootChildren(t *testing.T) {
	var (
		err     error
		laby    image.Image
		scanner binimg.Scanner
	)
	laby, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	scanner, err = binimg.NewScanner(laby)
	check(t, err)

	var testTbl = []struct {
		res   int  // resolution
		dir   side // direction
		white int  // num white children
		black int  // num black children
	}{
		{8, north, 2, 1},
		{8, south, 2, 1},
		{8, east, 1, 1},
		{8, west, 3, 1},
	}

	for _, tt := range testTbl {
		q, err := NewBUQuadtree(scanner, tt.res)
		check(t, err)

		root := q.root
		if root == nil {
			t.Fatalf("resolution %d, quadtree root is nil, expected not nil",
				tt.res)
		}

		var children QNodeList
		root.children(tt.dir, &children)
		var black, white int
		for _, nb := range children {
			switch nb.Color() {
			case Black:
				black++
			case White:
				white++
			}
		}
		if tt.white != white {
			t.Errorf("resolution %d, expected root to have %d white children at %s, got %d",
				tt.res, tt.white, tt.dir, white)

		}
		if tt.black != black {
			t.Errorf("resolution %d, expected root to have %d black children at %s, got %d",
				tt.res, tt.black, tt.dir, black)
		}
	}
}

func TestBUQuadtreeChildren(t *testing.T) {
	var (
		err     error
		laby    image.Image
		scanner binimg.Scanner
	)
	laby, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	scanner, err = binimg.NewScanner(laby)
	check(t, err)

	var testTbl = []struct {
		res   int         // resolution
		pt    image.Point // queried point child
		dir   side        // direction
		white int         // num white children
		black int         // num black children
	}{
		{8, image.Pt(8, 8), north, 2, 0},
		{8, image.Pt(8, 8), south, 1, 1},
		{8, image.Pt(8, 8), east, 2, 0},
		{8, image.Pt(8, 8), west, 1, 1},
	}

	for _, tt := range testTbl {
		q, err := NewBUQuadtree(scanner, tt.res)
		check(t, err)

		node := q.PointLocation(tt.pt)
		exists := node != nil
		if !exists {
			t.Fatalf("resolution %d, expected exists to be true for point %v, got false instead",
				tt.res, tt.pt)
		}

		parent := node.Parent().(*BUQNode)
		if parent == nil {
			t.Fatalf("resolution %d, parent of %v is nil, expected not nil",
				tt.res, tt.pt)
		}

		var children QNodeList
		parent.children(tt.dir, &children)
		var black, white int
		for _, nb := range children {
			switch nb.Color() {
			case Black:
				black++
			case White:
				white++
			}
		}
		if tt.white != white {
			t.Errorf("resolution %d, expected %v to have %d white children at %s, got %d",
				tt.res, parent, tt.white, tt.dir, white)

		}
		if tt.black != black {
			t.Errorf("resolution %d, expected %v to have %d black children at %s, got %d",
				tt.res, parent, tt.black, tt.dir, black)
		}
	}
}
