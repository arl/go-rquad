package rquad

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/go-rquad/internal"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestBasicBinImageTree(t *testing.T) {
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
		scanner imgscan.Scanner
	)

	for _, tt := range testTbl {
		bm, err = internal.LoadPNG(tt.fn)
		check(t, err)
		scanner, err = imgscan.NewScanner(bm)
		check(t, err)

		for _, res := range tt.resolutions {
			model, err := NewBinImgTreeModel(scanner, res)
			check(t, err)
			q := NewBasicTree(model)

			var white, black int
			q.ForEachLeaf(func(n Node) {
				switch n.(*ColoredNode).color {
				case color.White:
					white++
				case color.Black:
					black++
				default:
					t.Fatalf("node color is neither white nor black, color: %v\n", n.(*ColoredNode).color)
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
