package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

type newQuadtreeFunc func(binimg.Scanner, int) (Quadtree, error)

func newBUQuadtree(scanner binimg.Scanner, resolution int) (Quadtree, error) {
	return NewBUQuadtree(scanner, resolution)
}

func testQuadtreeWhiteNodes(t *testing.T, fn newQuadtreeFunc) {
	var testTbl = []struct {
		fn          string // filename
		resolutions []int  // various resolutions
		expected    int    // number of expected white nodes
	}{
		{"./testdata/labyrinth1.32x32.png", []int{1, 2, 3, 4, 5, 6, 7, 8}, 7},
		{"./testdata/labyrinth1.32x32.png", []int{9, 15}, 1},
		{"./testdata/labyrinth2.32x32.png", []int{1, 2, 3, 4}, 33},
		{"./testdata/labyrinth2.32x32.png", []int{5, 6, 7, 8, 9, 15}, 0},
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

			whiteNodes := q.WhiteNodes()
			if len(whiteNodes) != tt.expected {
				t.Errorf("resolution:%d, expected %d white nodes, got %d",
					tt.expected, res, len(whiteNodes))
			}
		}

	}
}

func TestBUQuadtreeWhiteNodes(t *testing.T) {
	testQuadtreeWhiteNodes(t, newBUQuadtree)
}

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, resolution int) {
	var (
		bm      image.Image
		err     error
		scanner binimg.Scanner
	)

	bm, err = loadPNG(pngfile)
	checkB(b, err)
	scanner, err = binimg.NewScanner(bm)
	checkB(b, err)

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := NewBUQuadtree(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkBUQuadtreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 2)
}

func BenchmarkBUQuadtreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 4)
}

func BenchmarkBUQuadtreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 8)
}

func BenchmarkBUQuadtreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 16)
}
