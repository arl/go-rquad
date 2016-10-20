package rquad

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func TestBasicTreeLogicalErrors(t *testing.T) {
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
		_, err = NewBasicTree(scanner, tt.res)
		actual := err == nil
		if actual != tt.expected {
			t.Errorf("(%d,%d,%d): expected %v, actual %v, err:'%v'",
				tt.w, tt.h, tt.res, tt.expected, actual, err)
		}
	}
}

func TestBasicTreeQuery(t *testing.T) {
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
		q, err := NewBasicTree(scanner, res)
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
