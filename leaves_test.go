package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func testQuadtreeCountLeaves(t *testing.T, fn newQuadtreeFunc) {
	var testTbl = []struct {
		fn          string // filename
		resolutions []int  // various resolutions
		white       int    // number of expected white nodes
		black       int    // number of expected black nodes
	}{
		{"./testdata/labyrinth1.32x32.png", []int{1, 2, 3, 4, 5, 6, 7, 8}, 7, 3},
		{"./testdata/labyrinth1.32x32.png", []int{9, 15}, 1, 3},
		{"./testdata/labyrinth1.32x32.png", []int{16}, 1, 3},
		{"./testdata/labyrinth2.32x32.png", []int{1, 2, 3, 4}, 33, 19},
		{"./testdata/labyrinth2.32x32.png", []int{5, 6, 7, 8}, 0, 16},
		{"./testdata/labyrinth2.32x32.png", []int{9, 15}, 0, 4},
		{"./testdata/labyrinth3.32x32.png", []int{1, 2, 3, 4}, 7, 6},
		{"./testdata/labyrinth3.32x32.png", []int{5, 6, 7, 8}, 5, 5},
		{"./testdata/labyrinth3.32x32.png", []int{9, 15}, 1, 3},
	}
	var (
		err     error
		bm      image.Image
		scanner binimg.Scanner
	)

	for _, tt := range testTbl {
		bm, err = loadPNG(tt.fn)
		check(t, err)
		scanner, err = binimg.NewScanner(bm)
		check(t, err)

		for _, res := range tt.resolutions {
			q, err := fn(scanner, res)
			check(t, err)

			var white, black int
			q.ForEachLeaf(Gray, func(n Node) {
				switch n.Color() {
				case White:
					white++
				case Black:
					black++
				case Gray:
					t.Fatalf("got gray leaf node")
				}
			})
			if white != tt.white {
				t.Errorf("on %s resolution:%d, expected %d white nodes, got %d",
					tt.fn, res, tt.white, white)
			}
			if black != tt.black {
				t.Errorf("on %s resolution:%d, expected %d black nodes, got %d",
					tt.fn, res, tt.black, black)
			}
		}
	}
}

func TestBUQuadtreeCountLeaves(t *testing.T) {
	testQuadtreeCountLeaves(t, newBUQuadtree)
}

func TestCNQuadtreeCountLeaves(t *testing.T) {
	testQuadtreeCountLeaves(t, newCNQuadtree)
}
