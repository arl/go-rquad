package quadtree

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/binimg"
)

// BUQuadtree is a standard quadtree implementation with bottom-up neighbor
// finding technique.
//
// BUQuadtree works on rectangles quadrants as well as squares; quadrants of
// the same parent may have different dimensions due to the integer division.
// It internally handles BUQNode's which implement the QNode interface.
type BUQuadtree struct {
	resolution int            // maximal resolution
	scanner    binimg.Scanner // reference image
	root       *BUQNode       // root node
	leaves     QNodeList      // leaf nodes (filled during creation)
}

// NewBUQuadtree creates a BUQuadtree and populates it with BUQNode's,
// according to the content of the scanned image.
//
// resolution is the minimal dimension that can have a leaf node, no further
// subdivisions will be performed on a node if its width or height is equal to
// this value.
func NewBUQuadtree(scanner binimg.Scanner, resolution int) (*BUQuadtree, error) {
	// initialize package level variables
	initPackage()

	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}

	// To ensure a consistent behavior and eliminate corner cases,
	// the Quadtree's root node needs to have children. Thus, the
	// first instantiated BUQNode needs to always be subdivided.
	// This condition asserts the resolution is respected.
	minDim := scanner.Bounds().Dx()
	if scanner.Bounds().Dy() < minDim {
		minDim = scanner.Bounds().Dy()
	}
	if minDim < resolution*2 {
		return nil, errors.New("the image smaller dimension must be greater or equal to twice the resolution")
	}

	q := &BUQuadtree{
		resolution: resolution,
		scanner:    scanner,
	}

	q.root = q.createRootNode()
	return q, nil
}

// ForEachLeaf calls the given function for each leaf node of the quadtree.
//
// Successive calls to the provided function are performed in no particular
// order. The color parameter allows to loop on the leaves of a particular
// color, Black or White.
// NOTE: As by definition, Gray leaves do not exist, passing Gray to
// ForEachLeaf should return all leaves, independently of their color.
func (q *BUQuadtree) ForEachLeaf(color QNodeColor, fn func(QNode)) {
	for _, n := range q.leaves {
		if color == Gray || n.Color() == color {
			fn(n)
		}
	}
}

func (q *BUQuadtree) createRootNode() *BUQNode {
	n := &BUQNode{
		color:  Gray,
		bounds: q.scanner.Bounds(),
	}
	q.subdivide(n)
	return n
}

func (q *BUQuadtree) createInnerNode(bounds image.Rectangle, parent *BUQNode, location quadrant) *BUQNode {
	n := &BUQNode{
		color:    Gray,
		bounds:   bounds,
		parent:   parent,
		location: location,
	}

	uniform, col := q.scanner.Uniform(bounds)
	switch uniform {
	case true:
		// quadrant is uniform, won't need to subdivide any further
		if col == binimg.White {
			n.color = White
		} else {
			n.color = Black
		}
	case false:
		// if we reached maximal resolution..
		if n.bounds.Dx()/2 < q.resolution || n.bounds.Dy()/2 < q.resolution {
			// ...make this node a black leaf, instead of gray
			n.color = Black
		} else {
			q.subdivide(n)
		}
	}

	// fills leaves slices
	if n.color != Gray {
		q.leaves = append(q.leaves, n)
	}
	return n
}

func (q *BUQuadtree) subdivide(n *BUQNode) {
	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'
	//
	x0 := n.bounds.Min.X
	x1 := n.bounds.Min.X + n.bounds.Dx()/2
	x2 := n.bounds.Max.X

	y0 := n.bounds.Min.Y
	y1 := n.bounds.Min.Y + n.bounds.Dy()/2
	y2 := n.bounds.Max.Y

	// create the 4 children nodes, one per quadrant
	n.northWest = q.createInnerNode(image.Rect(x0, y0, x1, y1), n, northWest)
	n.southWest = q.createInnerNode(image.Rect(x0, y1, x1, y2), n, southWest)
	n.northEast = q.createInnerNode(image.Rect(x1, y0, x2, y1), n, northEast)
	n.southEast = q.createInnerNode(image.Rect(x1, y1, x2, y2), n, southEast)
}

// Root returns the quadtree root node.
func (q *BUQuadtree) Root() QNode {
	return q.root
}
