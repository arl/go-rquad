package rquad

import (
	"image"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func testQuadtreeNeighbours(t *testing.T, fn newQuadtreeFunc) {
	var (
		laby1, laby2 *binimg.Binary
		err          error
	)

	// load both test images
	laby1, err = loadPNG("./testdata/labyrinth1.32x32.png")
	check(t, err)
	laby2, err = loadPNG("./testdata/labyrinth4.8x8.png")
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
		{laby1, 8, image.Pt(3, 3), 1, 1},
		{laby1, 8, image.Pt(11, 3), 2, 1},
		{laby1, 8, image.Pt(23, 7), 3, 0},
		{laby1, 8, image.Pt(3, 11), 3, 0},
		{laby1, 8, image.Pt(11, 11), 2, 2},
		{laby1, 8, image.Pt(3, 19), 2, 1},
		{laby1, 8, image.Pt(11, 19), 3, 1},
		{laby1, 8, image.Pt(23, 23), 1, 2},
		{laby1, 8, image.Pt(3, 27), 1, 1},
		{laby1, 8, image.Pt(11, 27), 3, 0},
		{laby1, 16, image.Pt(11, 27), 1, 1},

		{laby2, 1, image.Pt(0, 0), 1, 1},
		{laby2, 1, image.Pt(2, 0), 2, 2},
		{laby2, 1, image.Pt(4, 0), 2, 1},
		{laby2, 1, image.Pt(2, 2), 2, 2},
		{laby2, 1, image.Pt(3, 2), 3, 1},
		{laby2, 1, image.Pt(2, 3), 4, 0},
		{laby2, 1, image.Pt(3, 3), 3, 1},
		{laby2, 1, image.Pt(6, 0), 1, 1},
		{laby2, 1, image.Pt(4, 2), 3, 2},
		{laby2, 1, image.Pt(0, 4), 4, 0},
	}

	for _, tt := range testTbl {
		scanner, err := binimg.NewScanner(tt.img)
		check(t, err)
		q, err := fn(scanner, tt.res)
		check(t, err)

		node := PointLocation(q, tt.pt)
		exists := node != nil
		if !exists {
			t.Fatalf("%s, resolution %d, expected exists to be true for point %v, got false instead",
				imgAlias[tt.img], tt.res, tt.pt)
		}

		white, black := neighbourColors(node)
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

func TestBasicQuadtreeNeighbours(t *testing.T) {
	testQuadtreeNeighbours(t, newBasicTree)
}

func TestCNTreeNeighbours(t *testing.T) {
	testQuadtreeNeighbours(t, newCNTree)
}

// TODO: this test should, obviously, test the results obtained to the correct
// ones, instead of comparing the results obtained with two different quadtree
// implementations. At least, like that, we know the results are homogeneous, we
// are just not sure they are correct! ;-)
func TestNeighboursFinding(t *testing.T) {
	var (
		img     *binimg.Binary
		scanner binimg.Scanner
		err     error
	)
	img, err = loadPNG("./testdata/bigsquare.png")
	check(t, err)

	r := rand.New(rand.NewSource(99))

	scanner, err = binimg.NewScanner(img)
	check(t, err)

	// create a cardinal neighbour and a basic quadtree
	card, err := NewCNTree(scanner, 8)
	check(t, err)
	basic, err := NewBasicTree(scanner, 8)
	check(t, err)

	// TODO: check with real image bounds
	if card.Root().Bounds() != basic.Root().Bounds() {
		t.Fatalf("got different bounds, wanted equal")
	}

	randomPt := func(rect image.Rectangle) image.Point {
		return image.Pt(r.Intn(rect.Max.X-rect.Min.X)+rect.Min.X,
			r.Intn(rect.Max.Y-rect.Min.Y)+rect.Min.Y)
	}

	for i := 0; i < 50; i++ {
		pt := randomPt(card.Root().Bounds())

		cnnode := PointLocation(card, pt)
		basicnode := PointLocation(basic, pt)
		if (cnnode != nil) != (basicnode != nil) {
			t.Errorf("got different node existence for point %v, wanted the same", pt)
		}

		cnwhite, cnblack := neighbourColors(cnnode)
		bwhite, bblack := neighbourColors(basicnode)
		if cnwhite != bwhite {
			t.Errorf("got %d white neighbours for cnnode, %d for basicnode, wanted the same number", cnwhite, bwhite)
		}
		if cnblack != bblack {
			t.Errorf("got %d black neighbours for cnnode, %d for basicnode, wanted the same number", cnblack, bblack)
		}
	}
}

func benchmarkNeighboursFinding(b *testing.B, fn newQuadtreeFunc, numPoints int) {
	var (
		img     *binimg.Binary
		scanner binimg.Scanner
		err     error
	)
	img, err = loadPNG("./testdata/bigsquare.png")
	checkB(b, err)

	r := rand.New(rand.NewSource(99))

	scanner, err = binimg.NewScanner(img)
	checkB(b, err)

	// create a cardinal neighbour and a basic quadtree
	q, err := fn(scanner, 8)
	checkB(b, err)

	randomPt := func(rect image.Rectangle) image.Point {
		return image.Pt(r.Intn(rect.Max.X-rect.Min.X)+rect.Min.X,
			r.Intn(rect.Max.Y-rect.Min.Y)+rect.Min.Y)
	}

	noop := func(Node) {}

	// fill a slice with random points
	points := make([]image.Point, numPoints, numPoints)
	for i := 0; i < numPoints; i++ {
		pt := randomPt(q.Root().Bounds())
		points[i] = pt
	}

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, pt := range points {
			node := PointLocation(q, pt)
			ForEachNeighbour(node, noop)
		}
	}
}

func BenchmarkBasicQuadtreeNeighboursFinding10(b *testing.B) {
	benchmarkNeighboursFinding(b, newBasicTree, 10)
}

func BenchmarkBasicQuadtreeNeighboursFinding50(b *testing.B) {
	benchmarkNeighboursFinding(b, newBasicTree, 50)
}

func BenchmarkBasicQuadtreeNeighboursFinding200(b *testing.B) {
	benchmarkNeighboursFinding(b, newBasicTree, 200)
}

func BenchmarkBasicQuadtreeNeighboursFinding1000(b *testing.B) {
	benchmarkNeighboursFinding(b, newBasicTree, 1000)
}

func BenchmarkCNTreeNeighboursFinding10(b *testing.B) {
	benchmarkNeighboursFinding(b, newCNTree, 10)
}

func BenchmarkCNTreeNeighboursFinding50(b *testing.B) {
	benchmarkNeighboursFinding(b, newCNTree, 50)
}

func BenchmarkCNTreeNeighboursFinding200(b *testing.B) {
	benchmarkNeighboursFinding(b, newCNTree, 200)
}

func BenchmarkCNTreeNeighboursFinding1000(b *testing.B) {
	benchmarkNeighboursFinding(b, newCNTree, 1000)
}
