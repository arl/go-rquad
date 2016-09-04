package quadtree

import (
	"testing"

	"github.com/RookieGameDevs/quadtree/bmp"
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

		_, err := NewQuadtreeFromBitmap(bm, tt.res)
		actual := err == nil
		if actual != tt.expected {
			t.Errorf("TestQuadtreeLogicalErrors (%d,%d,%d): expected %v, actual %v, err:'%v'", tt.w, tt.h, tt.res, tt.expected, actual, err)
		}
	}
}
