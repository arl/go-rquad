package quadtree

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

func TestQuadtreeLogicalErrors(t *testing.T) {
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

		_, err := NewQuadtreeFromBitmap(bm, tt.res, nil)
		actual := err == nil
		if actual != tt.expected {
			t.Errorf("TestQuadtreeLogicalErrors (%d,%d,%d): expected %v, actual %v, err:'%v'", tt.w, tt.h, tt.res, tt.expected, actual, err)
		}
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuadtreeSubdivisions(t *testing.T) {
	nodes := []*quadnode{}
	saveNode := func(n *quadnode) {
		nodes = append(nodes, n)
	}

	// this is a simple 32x32 image, white background with 3 black squares the
	// biggest of which is 8x8 pixels, meaning that no nodes can ever be
	// smaller than 8x8, that's why every resolutions lower or equal than 8
	// should produce the same number of nodes.
	f, err := os.Open("./testdata/labyrinth.32x32.png")
	check(t, err)
	defer f.Close()

	var (
		img image.Image
		bm  *bmp.Bitmap
	)
	img, err = png.Decode(f)
	check(t, err)

	bm = bmp.NewFromImage(img)

	for _, res := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		_, err := NewQuadtreeFromBitmap(bm, res, saveNode)
		check(t, err)
		if len(nodes) != 13 {
			t.Errorf("TestQuadtreeSubdivisions (res:%d): expected 13 nodes, got %d", res, len(nodes))
		}
		// clear the list of nodes before re-using it
		nodes = nil
	}
}
