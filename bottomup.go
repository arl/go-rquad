package quadtree

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

// BUQuadtree is a standard quadtree implementation with bottom-up neighbor
// finding technique.
//
// BUQuadtree works on rectangles quadrants as well as squares; quadrants of
// the same parent may have different dimensions due to the integer division.
// It internally handles BUQuadnode, that implement the Quadnode interface.
type BUQuadtree struct {
	resolution int
	bm         *bmp.Bitmap
	root       *BUQuadnode
}

// NewBUQuadtree creates a BUQuadtree and populates it with BUQuadnode,
// following the given bitmap.
func NewBUQuadtree(bm *bmp.Bitmap, resolution int) (*BUQuadtree, error) {
	// initialize package level variables
	initPackage()

	// To ensure a consistent behavior and eliminate corner cases, the
	// Quadtree's root node need to have children, i.e. it can't
	// be a leaf node. Thus, the first instantiated BUQuadnode need to
	// always be subdivided. These two conditions make sure that
	// even with this subdivision the resolution will be respected.
	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}
	minDim := bm.Width
	if bm.Height < minDim {
		minDim = bm.Height
	}
	if minDim < resolution*2 {
		return nil, errors.New("the bitmap smaller dimension must be greater or equal to twice the resolution")
	}

	q := &BUQuadtree{
		resolution: resolution,
		bm:         bm,
	}
	q.root = q.createRootNode()
	return q, nil
}

func (q *BUQuadtree) createRootNode() *BUQuadnode {
	n := &BUQuadnode{
		quadnode: quadnode{
			color:       bmp.Gray,
			topLeft:     image.Point{0, 0},
			bottomRight: image.Point{q.bm.Width, q.bm.Height},
		},
	}
	q.subdivide(n)
	return n
}

func (q *BUQuadtree) createInnerNode(topleft, bottomright image.Point, parent *BUQuadnode) *BUQuadnode {
	n := &BUQuadnode{
		quadnode: quadnode{
			color:       bmp.Gray,
			topLeft:     topleft,
			bottomRight: bottomright,
			parent:      parent,
		},
	}

	n.color = q.bm.IsFilled(topleft, bottomright)

	switch {
	case n.width() <= q.resolution || n.height() <= q.resolution:
		// reached the maximal resolution
		break
	case n.color == bmp.Gray:
		q.subdivide(n)
	default:
		// quadrant is monocolor, don't need any further subdivisions
		break
	}
	return n
}

func (q *BUQuadtree) subdivide(n *BUQuadnode) {
	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'

	x0 := n.topLeft.X
	x1 := n.topLeft.X + n.width()/2
	x2 := n.bottomRight.X

	y0 := n.topLeft.Y
	y1 := n.topLeft.Y + n.height()/2
	y2 := n.bottomRight.Y

	// create the 4 children nodes, one per quadrant
	n.northWest = q.createInnerNode(
		image.Point{x0, y0},
		image.Point{x1, y1},
		n)
	n.southWest = q.createInnerNode(
		image.Point{x0, y1},
		image.Point{x1, y2},
		n)
	n.northEast = q.createInnerNode(
		image.Point{x1, y0},
		image.Point{x2, y1},
		n)
	n.southEast = q.createInnerNode(
		image.Point{x1, y1},
		image.Point{x2, y2},
		n)
}

// PointQuery returns the Quadnode containing the point at given coordinates.
//
// If such node doesn't exist, exists is false.
func (q *BUQuadtree) PointQuery(pt image.Point) (n Quadnode, exists bool) {
	return q.root.pointQuery(pt)
}
