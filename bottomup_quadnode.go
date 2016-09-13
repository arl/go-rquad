package quadtree

import (
	"image"
	"math"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

// BUQuadnode is a node of a BUQuadtree.
//
// It is a basic implementation of the Quadnode interface, augmented with
// methods implementing the bottom-up neighbor finding techniques.
type BUQuadnode struct {
	quadnode
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *BUQuadnode) isLeaf() bool {
	return n.color != bmp.Gray
}

// isLeaf checks if this node is white.
func (n *BUQuadnode) isWhite() bool {
	return n.color == bmp.White
}

// children fills the given slice with all the leaf children of this node (i.e
// either black or white), that can be found in a given direction.
func (n *BUQuadnode) children(dir side, nodes *[]*BUQuadnode) {

	if n.isLeaf() {
		return
	}

	var (
		s1, s2 *BUQuadnode
	)

	switch dir {
	case north:
		s1 = n.northEast.(*BUQuadnode)
		s2 = n.northWest.(*BUQuadnode)
		break
	case east:
		s1 = n.northEast.(*BUQuadnode)
		s2 = n.southEast.(*BUQuadnode)
		break
	case south:
		s1 = n.southEast.(*BUQuadnode)
		s2 = n.southWest.(*BUQuadnode)
		break
	case west:
		s1 = n.southEast.(*BUQuadnode)
		s2 = n.southWest.(*BUQuadnode)
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
func (n *BUQuadnode) equalSizeNeighbor(dir side) *BUQuadnode {
	var neighbor *BUQuadnode

	// Ascent the tree up to a common ancestor.
	if n.parent != nil {
		buparent := n.parent.(*BUQuadnode)
		if adjacent(dir, n.quadrant()) {
			neighbor = buparent.equalSizeNeighbor(dir)
		} else {
			neighbor = buparent
		}
	}

	// Backtrack mirroring the ascending moves.
	if neighbor != nil && !neighbor.isLeaf() {
		return neighbor.child(reflect(dir, n.quadrant())).(*BUQuadnode)
	} else {
		return neighbor
	}
}

// cornerNeighbor locates a neighbor of the current quadnode in the horizontal
// or vertical direction which is adjacent to one of its corners.
//
// The neighboring node must be adjacent to this corner.
// It can return nil if the neighbor can't be found.
func (n *BUQuadnode) cornerNeighbor(dir side, corner quadrant) *BUQuadnode {

	// If no neighbor can be found in the given
	// direction, node will be nil
	node := n.equalSizeNeighbor(dir)
	if node == nil {
		return nil
	}

	// Go down until we reach either a free or
	// an obstructed node, i.e. a leaf node.
	for !node.isLeaf() {
		node = node.child(reflect(dir, corner)).(*BUQuadnode)
	}

	return node
}

// _neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *BUQuadnode) _neighbours(dir side, nodes *[]*BUQuadnode) {

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
func (n *BUQuadnode) neighbours() []*BUQuadnode {

	nodes := make([]*BUQuadnode, 0)
	n._neighbours(north, &nodes)
	n._neighbours(south, &nodes)
	n._neighbours(east, &nodes)
	n._neighbours(west, &nodes)
	return nodes
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

func (n *BUQuadnode) pointQuery(pt image.Point) (Quadnode, bool) {
	if !n.inbound(pt) {
		return nil, false
	}
	if n.color != bmp.Gray {
		return n, true
	}
	nw := n.northWest.(*BUQuadnode)
	ne := n.northEast.(*BUQuadnode)
	sw := n.southWest.(*BUQuadnode)
	se := n.southEast.(*BUQuadnode)

	if nw.inbound(pt) {
		return nw.pointQuery(pt)
	} else if ne.inbound(pt) {
		return ne.pointQuery(pt)
	} else if sw.inbound(pt) {
		return sw.pointQuery(pt)
	} else if se.inbound(pt) {
		return se.pointQuery(pt)
	}
	panic("should never be here")
}
