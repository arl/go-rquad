package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
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

	for _, tt := range testTbl {
		bm := bmp.New(tt.w, tt.h)

		_, err := NewBUQuadtree(bm, tt.res)
		actual := err == nil
		if actual != tt.expected {
			t.Errorf("(%d,%d,%d): expected %v, actual %v, err:'%v'", tt.w, tt.h, tt.res, tt.expected, actual, err)
		}
	}
}

func TestBUQuadtreeSubdivisions(t *testing.T) {
	// this is a simple 32x32 image, white background with 3 black squares,
	// located so that they fill a quadrant the biggest of which is 8x8 pixels,
	// meaning that no nodes can ever be smaller than 8x8, that's why every
	// resolutions lower or equal than 8 should produce the same number of
	// nodes.
	bm, err := loadPNG("./testdata/labyrinth.32x32.png")
	check(t, err)

	for _, res := range []int{1, 2, 3, 4, 5, 6, 7, 8, 15} {
		q, err := NewBUQuadtree(bm, res)
		check(t, err)

		nodes := listNodes(q.root)
		if len(nodes) != 7 {
			t.Errorf("TestBUQuadtreeSubdivisions: resolution:%d, expected 7 nodes, got %d", res, len(nodes))
		}
	}
}

func TestBUQuadtreePointQuery(t *testing.T) {
	var testTbl = []struct {
		pt     image.Point // queried point
		exists bool        //node exists
		eqRef  bool        // node should be equal to ref node
	}{
		{image.Point{8, 0}, true, true},
		{image.Point{15, 0}, true, true},
		{image.Point{8, 7}, true, true},
		{image.Point{15, 7}, true, true},
		{image.Point{7, 0}, true, false},
		{image.Point{16, 0}, true, false},
		{image.Point{8, 8}, true, false},
		{image.Point{16, 8}, true, false},
		{image.Point{1, 31}, true, false},
		{image.Point{-1, 0}, false, false},
		{image.Point{32, 0}, false, false},
		{image.Point{0, -1}, false, false},
		{image.Point{0, 32}, false, false},
	}

	// coordinate of reference node
	refPt := image.Point{8, 0}

	bm, err := loadPNG("./testdata/labyrinth.32x32.png")
	check(t, err)

	for _, res := range []int{1, 2, 3, 4, 5, 6, 7, 8, 15} {
		q, err := NewBUQuadtree(bm, res)
		check(t, err)

		for _, tt := range testTbl {
			node, exists := q.PointQuery(tt.pt)
			if exists != tt.exists {
				t.Fatalf("resolution %d, expected exists to be %t for point %v, got %t instead", res, tt.exists, tt.pt, exists)
			}
			// obtain the refNode
			refNode, refExists := q.PointQuery(refPt)
			if !refExists {
				t.Fatalf("reference node should exist")
			}

			if exists {
				// check if queried node should be equal to reference node
				eqRef := node == refNode
				if eqRef != tt.eqRef {
					t.Errorf("resolution %d, got %t for comparison between reference node and queried node, expected %t. node:%v", res, eqRef, tt.eqRef, node)
				}
			}
		}
	}
}
