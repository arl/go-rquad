package quadtree

import (
	"errors"
	"image"

	"github.com/RookieGameDevs/quadtree/bmp"
)

// BUQuadtree is a standard quadtree implementation with bottom-up neighbor
// finding technique.
//
// It internally uses the BUQuadnode implementation of Quadnode.
type BUQuadtree struct {
	resolution int
	bm         *bmp.Bitmap
	root       *BUQuadnode
}

func NewBUQuadtree(bm *bmp.Bitmap, resolution int) (*BUQuadtree, error) {
	// To ensure a consistent behavior and eliminate corner cases, the
	// Quadtree's root node need to have children, i.e. it can't
	// be a leaf node. Thus, the first instantiated Quadnode need to
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
	q.root = q.CreateRootNode().(*BUQuadnode)
	return q, nil
}

func (q *BUQuadtree) CreateRootNode() Quadnode {
	n := &BUQuadnode{
		quadnode: quadnode{
			color:       bmp.Gray,
			topLeft:     image.Point{0, 0},
			bottomRight: image.Point{q.bm.Width, q.bm.Height},
		},
	}
	q.Subdivide(n)
	return n
}

func (q *BUQuadtree) CreateInnerNode(topleft, bottomright image.Point, parent Quadnode) Quadnode {
	n := &BUQuadnode{
		quadnode: quadnode{
			color:       bmp.Gray,
			topLeft:     topleft,
			bottomRight: bottomright,
			parent:      parent.(*BUQuadnode),
		},
	}

	n.color = q.bm.IsFilled(topleft, bottomright)

	switch {
	case n.width() <= q.resolution || n.height() <= q.resolution:
		// reached the maximal resolution
		break
	case n.color == bmp.Gray:
		q.Subdivide(n)
	default:
		// quadrant is monocolor, don't need any further subdivisions
		break
	}
	return n
}

func (q *BUQuadtree) Subdivide(n Quadnode) {
	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'

	node := n.(*BUQuadnode)

	x0 := node.topLeft.X
	x1 := node.topLeft.X + node.width()/2
	x2 := node.bottomRight.X

	y0 := node.topLeft.Y
	y1 := node.topLeft.Y + node.height()/2
	y2 := node.bottomRight.Y

	// create the 4 children nodes, one per quadrant
	node.northWest = q.CreateInnerNode(
		image.Point{x0, y0},
		image.Point{x1, y1},
		n).(*BUQuadnode)
	node.southWest = q.CreateInnerNode(
		image.Point{x0, y1},
		image.Point{x1, y2},
		n).(*BUQuadnode)
	node.northEast = q.CreateInnerNode(
		image.Point{x1, y0},
		image.Point{x2, y1},
		n).(*BUQuadnode)
	node.southEast = q.CreateInnerNode(
		image.Point{x1, y1},
		image.Point{x2, y2},
		n).(*BUQuadnode)
}
