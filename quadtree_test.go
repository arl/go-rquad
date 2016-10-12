package quadtree

import (
	"fmt"
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

type newQuadtreeFunc func(binimg.Scanner, int) (Quadtree, error)

func newBUQuadtree(scanner binimg.Scanner, resolution int) (Quadtree, error) {
	return NewBUQuadtree(scanner, resolution)
}

func newCNQuadtree(scanner binimg.Scanner, resolution int) (Quadtree, error) {
	return NewCNQuadtree(scanner, resolution)
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
		{"./testdata/labyrinth3.32x32.png", []int{1, 2, 3, 4}, 7},
		{"./testdata/labyrinth3.32x32.png", []int{5, 6, 7, 8}, 5},
		{"./testdata/labyrinth3.32x32.png", []int{9, 15}, 1},
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
				t.Errorf("on %s resolution:%d, expected %d white nodes, got %d",
					tt.fn, res, tt.expected, len(whiteNodes))
			}
		}
	}
}

func testQuadtreeNeighbours(t *testing.T, fn newQuadtreeFunc) {
	var (
		laby1, laby2 *binimg.Binary
		err          error
	)

	// load both test images
	laby1, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	laby2, err = loadPNG("./testdata/cn-8x8-3.png")
	check(t, err)

	// for logging purposes
	imgAlias := map[*binimg.Binary]string{
		laby1: "'labyrinth1.32x32'",
		laby2: "'cn-8x8-3.png'",
	}

	var testTbl = []struct {
		img   *binimg.Binary // source image
		res   int            // resolution
		pt    image.Point    // queried point
		white int            // num white neighbours
		black int            // num black neighbours
	}{
		{laby1, 8, image.Point{3, 3}, 1, 1},
		{laby1, 8, image.Point{11, 3}, 2, 1},
		{laby1, 8, image.Point{23, 7}, 3, 0},
		{laby1, 8, image.Point{3, 11}, 3, 0},
		{laby1, 8, image.Point{11, 11}, 2, 2},
		{laby1, 8, image.Point{3, 19}, 2, 1},
		{laby1, 8, image.Point{11, 19}, 3, 1},
		{laby1, 8, image.Point{23, 23}, 1, 2},
		{laby1, 8, image.Point{3, 27}, 1, 1},
		{laby1, 8, image.Point{11, 27}, 3, 0},
		{laby1, 16, image.Point{11, 27}, 1, 1},

		{laby2, 1, image.Point{0, 0}, 1, 1},
		{laby2, 1, image.Point{2, 0}, 2, 2},
		{laby2, 1, image.Point{4, 0}, 2, 1},
		{laby2, 1, image.Point{2, 2}, 2, 2},
		{laby2, 1, image.Point{3, 2}, 3, 1},
		{laby2, 1, image.Point{2, 3}, 4, 0},
		{laby2, 1, image.Point{3, 3}, 3, 1},
		{laby2, 1, image.Point{6, 0}, 1, 1},
		{laby2, 1, image.Point{4, 2}, 3, 2},
		{laby2, 1, image.Point{0, 4}, 4, 0},
	}

	for _, tt := range testTbl {
		scanner, err := binimg.NewScanner(tt.img)
		check(t, err)
		q, err := fn(scanner, tt.res)
		check(t, err)

		node, exists := Query(q, tt.pt)
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				imgAlias[tt.img], tt.res, tt.pt)
		}

		var (
			black, white int
			nodes        QNodeList
			strW, strB   string
		)
		node.Neighbours(&nodes)
		for _, nb := range nodes {
			switch nb.Color() {
			case Black:
				strB += fmt.Sprintln(nb)
				black++
			case White:
				strW += fmt.Sprintln(nb)
				white++
			}
		}

		if tt.white != white {
			t.Errorf("%s, resolution %d, expected pt %v to have %d white neighbours, got %d",
				imgAlias[tt.img], tt.res, tt.pt, tt.white, white)
			t.FailNow()
		}
		if tt.black != black {
			t.Errorf("%s, resolution %d, expected pt %v to have %d black neighbours, got %d",
				imgAlias[tt.img], tt.res, tt.pt, tt.black, black)
			t.FailNow()
		}
	}
}

func TestBUQuadtreeWhiteNodes(t *testing.T) {
	testQuadtreeWhiteNodes(t, newBUQuadtree)
}

func TestCNQuadtreeWhiteNodes(t *testing.T) {
	testQuadtreeWhiteNodes(t, newCNQuadtree)
}

func TestBUQuadtreeNeighbours(t *testing.T) {
	testQuadtreeNeighbours(t, newBUQuadtree)
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
		{image.Point{4, 2}, 3, 1},
		{image.Point{4, 4}, 3, 2},
		{image.Point{5, 3}, 3, 1},
		{image.Point{6, 2}, 3, 1},
	}

	scanner, err := binimg.NewScanner(laby)
	check(t, err)
	q, err := fn(scanner, 1)
	check(t, err)

	for _, tt := range testTbl {
		node, exists := Query(q, tt.pt)
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				pngfile, 1, tt.pt)
		}

		var black, white int
		var nodes QNodeList
		node.Neighbours(&nodes)
		var strW, strB string
		for _, nb := range nodes {
			switch nb.Color() {
			case Black:
				strB += fmt.Sprintln(nb.(*CNQNode))
				black++
			case White:
				strW += fmt.Sprintln(nb.(*CNQNode))
				white++
			}
		}
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
		{1, image.Point{0, 0}, 2, 0},
		{1, image.Point{1, 0}, 2, 1},
		{1, image.Point{2, 0}, 2, 1},
		{1, image.Point{3, 0}, 2, 0},
		{1, image.Point{0, 1}, 2, 1},
		{1, image.Point{1, 1}, 2, 2},
		{1, image.Point{2, 1}, 2, 2},
		{1, image.Point{3, 1}, 2, 1},
	}

	for _, tt := range testTbl {
		scanner, err := binimg.NewScanner(laby)
		check(t, err)
		q, err := fn(scanner, tt.res)
		check(t, err)

		node, exists := Query(q, tt.pt)
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				pngfile, tt.res, tt.pt)
		}

		var black, white int
		var nodes QNodeList
		node.Neighbours(&nodes)
		for _, nb := range nodes {
			switch nb.Color() {
			case Black:
				black++
			case White:
				white++
			}
		}
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

func BenchmarkBUQuadtreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBUQuadtree, 2)
}

func BenchmarkBUQuadtreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBUQuadtree, 4)
}

func BenchmarkBUQuadtreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBUQuadtree, 8)
}

func BenchmarkBUQuadtreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBUQuadtree, 16)
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
