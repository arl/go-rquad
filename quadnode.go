package quadtree

import (
	"image"

	"github.com/RookieGameDevs/quadtree/bmp"
)

//go:generate stringer -type=quadrant
type quadrant int

const (
	northWest quadrant = iota
	northEast
	southWest
	southEast
)

//go:generate stringer -type=side
type side int

const (
	north side = iota
	east
	south
	west
)

type quadnode struct {
	northWest *quadnode // pointer to the northwest child
	southWest *quadnode // pointer to the southwest child
	northEast *quadnode // pointer to the northeast child
	southEast *quadnode // pointer to the southeast child

	// node top-left corner coordinates, the origin
	topLeft image.Point

	// node bottom-right corner coordinates, the point is included
	bottomRight image.Point

	// node color
	color bmp.Color
}

func (n *quadnode) width() int {
	return n.bottomRight.X - n.topLeft.X
}

func (n *quadnode) height() int {
	return n.bottomRight.Y - n.topLeft.Y
}

func newRootQuadNode(bm *bmp.Bitmap, resolution int) *quadnode {
	n := &quadnode{
		color:       bmp.Gray,
		topLeft:     image.Point{0, 0},
		bottomRight: image.Point{bm.Width - 1, bm.Height - 1},
	}
	n.subdivide(bm, resolution)
	return n
}

// newInnerQuadNode construct a child node.
func newInnerQuadNode(bm *bmp.Bitmap, topLeft, bottomRight image.Point, resolution int, parent *quadnode) *quadnode {
	n := &quadnode{
		color:       bmp.Gray,
		topLeft:     topLeft,
		bottomRight: bottomRight,
	}

	color := bm.IsFilled(topLeft, bottomRight)
	switch {
	case n.width() <= resolution || n.height() <= resolution:
		fallthrough
	case color == bmp.Black:
		// quadrant is totally obstructed or we reached the maximal division,
		// no need to go further
		n.color = bmp.Black
	case color == bmp.White:
		// quadrant is totally empty, no need to go further
		n.color = bmp.White
	case color == bmp.Gray:
		n.subdivide(bm, resolution)
	}
	return n
}

// subdivide subdivides the current node into four children.
//
// This methode should be called once by the constructor if
// the current node intersect with an obstacle and its
// width and height are both greater than the resolution.
func (n *quadnode) subdivide(bm *bmp.Bitmap, resolution int) {
	// the default is gray
	n.color = bmp.Gray

	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'
	//     x0   x1     x2
	x0 := n.topLeft.X
	x1 := n.topLeft.X + n.width()/2
	x2 := n.bottomRight.X

	y0 := n.topLeft.Y
	y1 := n.topLeft.Y + n.height()/2
	y2 := n.bottomRight.Y

	// create the 4 children nodes, one per quadrant
	n.northWest = newInnerQuadNode(bm,
		image.Point{x0, y0},
		image.Point{x1, y1}, resolution, n)
	n.southWest = newInnerQuadNode(bm,
		image.Point{x0, y1},
		image.Point{x1, y2}, resolution, n)
	n.northEast = newInnerQuadNode(bm,
		image.Point{x1, y0},
		image.Point{x2, y1}, resolution, n)
	n.southEast = newInnerQuadNode(bm,
		image.Point{x1, y1},
		image.Point{x2, y2}, resolution, n)
}
