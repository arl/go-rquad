package quadtree

import (
	"image"
	"math"

	"github.com/RookieGameDevs/quadtree/bmp"
)

type quadnode struct {
	parent *quadnode // pointer to the parent node

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

func newRootQuadNode(bm *bmp.Bitmap, resolution int, fn PostNodeCreationFunc) *quadnode {
	n := &quadnode{
		color:       bmp.Gray,
		topLeft:     image.Point{0, 0},
		bottomRight: image.Point{bm.Width, bm.Height},
	}
	fn(n)
	n.subdivide(bm, resolution, fn)
	return n
}

// newInnerQuadNode construct a child node.
func newInnerQuadNode(bm *bmp.Bitmap, topLeft, bottomRight image.Point, resolution int, parent *quadnode, fn PostNodeCreationFunc) *quadnode {
	n := &quadnode{
		color:       bmp.Gray,
		topLeft:     topLeft,
		bottomRight: bottomRight,
		parent:      parent,
	}

	n.color = bm.IsFilled(topLeft, bottomRight)
	fn(n)

	switch {
	case n.width() <= resolution || n.height() <= resolution:
		// reached the maximal resolution
		break
	case n.color == bmp.Gray:
		n.subdivide(bm, resolution, fn)
	default:
		// quadrant is monocolor, don't need any further subdivisions
		break
	}
	return n
}

// subdivide subdivides the current node into four children.
//
// This methode should be called once by the constructor if
// the current node intersect with an obstacle and its
// width and height are both greater than the resolution.
func (n *quadnode) subdivide(bm *bmp.Bitmap, resolution int, fn PostNodeCreationFunc) {
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
	n.northWest = newInnerQuadNode(bm,
		image.Point{x0, y0},
		image.Point{x1, y1}, resolution, n, fn)
	n.southWest = newInnerQuadNode(bm,
		image.Point{x0, y1},
		image.Point{x1, y2}, resolution, n, fn)
	n.northEast = newInnerQuadNode(bm,
		image.Point{x1, y0},
		image.Point{x2, y1}, resolution, n, fn)
	n.southEast = newInnerQuadNode(bm,
		image.Point{x1, y1},
		image.Point{x2, y2}, resolution, n, fn)
}

// quadrant obtain this node's quadrant relative to its parent.
//
// must not be called on the root node
func (n *quadnode) quadrant() quadrant {
	if n.parent == nil {
		panic("the root node's quadrant is undefined")
	}

	if n.parent.northWest == n {
		return northWest
	} else if n.parent.southWest == n {
		return southWest
	} else if n.parent.northEast == n {
		return northEast
	} else {
		return southEast
	}
}

// child returns a pointer to the child node associated to the given quadrant
func (n *quadnode) child(q quadrant) *quadnode {
	switch q {
	case northWest:
		return n.northWest
	case southWest:
		return n.southWest
	case northEast:
		return n.northEast
	case southEast:
		return n.southEast
	}
	panic("undefined quadrant")
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *quadnode) isLeaf() bool {
	return n.color != bmp.Gray
}

// isLeaf checks if this node is white.
func (n *quadnode) isWhite() bool {
	return n.color == bmp.White
}

// children fills the given slice with all the leaf children of this node (i.e
// either black or white), that can be found in a given direction.
func (n *quadnode) children(dir side, nodes *[]*quadnode) {

	if n.isLeaf() {
		return
	}

	var (
		s1, s2 *quadnode
	)

	switch dir {
	case north:
		s1 = n.northEast
		s2 = n.northWest
		break
	case east:
		s1 = n.northEast
		s2 = n.southEast
		break
	case south:
		s1 = n.southEast
		s2 = n.southWest
		break
	case west:
		s1 = n.northEast
		s2 = n.southWest
	}

	if s1.isLeaf() {
		*nodes = append(*nodes, s1)
	} else {
		s1.children(dir, nodes)
	}

	if s2.isLeaf() {
		*nodes = append(*nodes, s2)
	} else {
		s2.children(dir, nodes)
	}
}

// equalSizeNeighbor locates an equal-sized neighbor of the current node in the
// vertical or horizontal direction.
//
//  cf. Hanan Samet 1981 article Neighbor Finding in Quadtrees.
// It can return nil if the neighbor can't be found.
func (n *quadnode) equalSizeNeighbor(dir side) *quadnode {
	var neighbor *quadnode

	// Ascent the tree up to a common ancestor.
	if n.parent != nil && adjacent(dir, n.quadrant()) {
		neighbor = n.parent.equalSizeNeighbor(dir)
	} else {
		neighbor = n.parent
	}

	// Backtrack mirroring the ascending moves.
	if neighbor != nil && !neighbor.isLeaf() {
		return neighbor.child(reflect(dir, n.quadrant()))
	} else {
		return neighbor
	}
}

// cornerNeighbor locates a neighbor of the current quadnode in the horizontal
// or vertical direction which is adjacent to one of its corners.
//
// The neighboring node must be adjacent to this corner.
// It can return nil if the neighbor can't be found.
func (n *quadnode) cornerNeighbor(dir side, corner quadrant) *quadnode {

	// If no neighbor can be found in the given
	// direction, node will be nil
	node := n.equalSizeNeighbor(dir)
	if node == nil {
		return nil
	}

	// Go down until we reach either a free or
	// an obstructed node, i.e. a leaf node.
	for !node.isLeaf() {
		node = node.child(reflect(dir, corner))
	}

	return node
}

// _neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *quadnode) _neighbours(dir side, nodes *[]*quadnode) {

	// If no neighbor can be found in the given
	// direction, node will be null.
	node := n.equalSizeNeighbor(dir)
	if node != nil {
		if node.isLeaf() {
			// Neighbor is already a leaf node, we're done.
			*nodes = append(*nodes, node)
		} else {
			// The neighbor isn't a leaf node so we need to
			// go further down matching its children, but in
			// the opposite direction from where we came.
			node.children(opposite(dir), nodes)
		}
	}
}

// neighbours returns a slice of all leaf neighbours of the current node.
func (n *quadnode) neighbours() []*quadnode {

	nodes := make([]*quadnode, 0)
	n._neighbours(north, &nodes)
	n._neighbours(south, &nodes)
	n._neighbours(east, &nodes)
	n._neighbours(west, &nodes)
	return nodes
}

func (n *quadnode) origin() image.Point {
	return n.topLeft
}

// squaredDistance returns the squared straight-line distance between this node and another.
func (n *quadnode) squaredDistance(other *quadnode) float64 {
	//         a    (x1,y1)
	//    .-----------.
	//    |        .-'
	//  b |     .-'
	//    |  .-'
	//    '-'
	// (x2, y2)

	x1 := float64(n.topLeft.X) + float64(n.width())/2
	y1 := float64(n.topLeft.Y) + float64(n.height())/2

	x2 := float64(other.topLeft.X) + float64(other.width())/2
	y2 := float64(other.topLeft.Y) + float64(other.height())/2

	a := math.Abs(x1 - x2)
	b := math.Abs(y1 - y2)

	return a*a + b*b
}

// inbound checks if a given point is inside the region represented by this
// node.
func (n *quadnode) inbound(pt image.Point) bool {
	return (n.topLeft.X <= pt.X && pt.X <= n.bottomRight.X) &&
		(n.topLeft.Y <= pt.Y && pt.Y <= n.bottomRight.Y)
}
