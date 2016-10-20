package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func TestBasicTreeNeighbours(t *testing.T) {
	testQuadtreeNeighbours(t, newBasicTree)
}

func TestCNQuadtreeNeighbours(t *testing.T) {
	testQuadtreeNeighbours(t, newCNQuadtree)
}

func testDebugQuadtreeNeighboursExample(t *testing.T, fn newQuadtreeFunc) {
	var (
		laby    *binimg.Binary
		err     error
		pngfile string
	)
	pngfile = "./testdata/labyrinth-small.8x8.step1.png"
	laby, err = loadPNG(pngfile)
	check(t, err)

	var testTbl = []struct {
		pt    image.Point // queried point
		white int         // num white neighbours
		black int         // num black neighbours
	}{
		{image.Pt(4, 2), 3, 1},
		{image.Pt(4, 4), 3, 2},
		{image.Pt(5, 3), 3, 1},
		{image.Pt(6, 2), 3, 1},
	}

	scanner, err := binimg.NewScanner(laby)
	check(t, err)
	q, err := fn(scanner, 1)
	check(t, err)

	for _, tt := range testTbl {
		node := q.(PointLocator).PointLocation(tt.pt)
		exists := node != nil
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				pngfile, 1, tt.pt)
		}

		white, black := neighbourColors(node)
		if tt.white != white {

			t.Errorf("%s, resolution %d, expected pt %v to have %d white neighbours, got %d",
				pngfile, 1, tt.pt, tt.white, white)
		}
		if tt.black != black {
			t.Errorf("%s, resolution %d, expected pt %v to have %d black neighbours, got %d",
				pngfile, 1, tt.pt, tt.black, black)
		}
	}
}

func testDebugQuadtreeNeighboursSmall(t *testing.T, fn newQuadtreeFunc) {
	var (
		laby    *binimg.Binary
		err     error
		pngfile string
	)
	pngfile = "./testdata/labyrinth-small.4x4.png"
	laby, err = loadPNG(pngfile)
	check(t, err)

	var testTbl = []struct {
		res   int         // resolution
		pt    image.Point // queried point
		white int         // num white neighbours
		black int         // num black neighbours
	}{
		{1, image.Pt(0, 0), 2, 0},
		{1, image.Pt(1, 0), 2, 1},
		{1, image.Pt(2, 0), 2, 1},
		{1, image.Pt(3, 0), 2, 0},
		{1, image.Pt(0, 1), 2, 1},
		{1, image.Pt(1, 1), 2, 2},
		{1, image.Pt(2, 1), 2, 2},
		{1, image.Pt(3, 1), 2, 1},
	}

	for _, tt := range testTbl {
		scanner, err := binimg.NewScanner(laby)
		check(t, err)
		q, err := fn(scanner, tt.res)
		check(t, err)

		node := q.(PointLocator).PointLocation(tt.pt)
		exists := node != nil
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				pngfile, tt.res, tt.pt)
		}

		white, black := neighbourColors(node)
		if tt.white != white {
			t.Errorf("%s, resolution %d, expected pt %v to have %d white neighbours, got %d",
				pngfile, tt.res, tt.pt, tt.white, white)
		}
		if tt.black != black {
			t.Errorf("%s, resolution %d, expected pt %v to have %d black neighbours, got %d",
				pngfile, tt.res, tt.pt, tt.black, black)
		}
	}
}

func TestDebugCNQuadtreeNeighbours(t *testing.T) {
	testDebugQuadtreeNeighboursExample(t, newCNQuadtree)
}

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, fn newQuadtreeFunc, resolution int) {
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
		_, err := fn(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkBasicTreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 2)
}

func BenchmarkBasicTreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 4)
}

func BenchmarkBasicTreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 8)
}

func BenchmarkBasicTreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 16)
}

func BenchmarkCNQuadtreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 2)
}

func BenchmarkCNQuadtreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 4)
}

func BenchmarkCNQuadtreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 8)
}

func BenchmarkCNQuadtreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 16)
}
