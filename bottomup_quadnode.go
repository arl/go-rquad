package quadtree

import "image"

// BUQuadnode is a node of a BUQuadtree.
//
// It is a basic implementation of the Quadnode interface, augmented with
// methods implementing the bottom-up neighbor finding techniques.
type BUQuadnode struct {
	quadnode
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *BUQuadnode) isLeaf() bool {
	return n.color != Gray
}

// children fills the given slice with all the leaf children of this node (i.e
// either black or white), that can be found in a given direction.
func (n *BUQuadnode) children(dir side, nodes *NodeList) {

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
		s1 = n.northWest.(*BUQuadnode)
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
	}
	return neighbor
}

// _neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *BUQuadnode) _neighbours(dir side, nodes *NodeList) {

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

// Neighbours fills a NodeList with the neighbours of this node. n must be
// a leaf node, or nodes will be an empty slice.
func (n *BUQuadnode) Neighbours(nodes *NodeList) {
	n._neighbours(north, nodes)
	n._neighbours(south, nodes)
	n._neighbours(east, nodes)
	n._neighbours(west, nodes)
}

// quadrant obtains this node's quadrant relative to its parent.
//
// must not be called on the root node
func (n *BUQuadnode) quadrant() quadrant {
	if n.parent.NorthWest() == n {
		return northWest
	} else if n.parent.SouthWest() == n {
		return southWest
	} else if n.parent.NorthEast() == n {
		return northEast
	}
	return southEast
}

func (n *BUQuadnode) pointQuery(pt image.Point) (Quadnode, bool) {
	if !n.inbound(pt) {
		return nil, false
	}
	if n.color != Gray {
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
	}
	return se.pointQuery(pt)
}
