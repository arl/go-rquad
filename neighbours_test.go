package quadtree

import (
	"image"
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
	laby2, err = loadPNG("./testdata/labyrinth4-8x8.png")
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

		node := q.(PointLocator).PointLocation(tt.pt)
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
