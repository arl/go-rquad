package quadtree

// BUQNode is a node of a BUQuadtree.
//
// It is a basic implementation of the QNode interface, augmented with
// methods implementing the bottom-up neighbor finding techniques.
type BUQNode struct {
	qnode
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *BUQNode) isLeaf() bool {
	return n.color != Gray
}

// children fills the given slice with all the leaf children of this node (i.e
// either black or white), that can be found in a given direction.
func (n *BUQNode) children(dir side, nodes *QNodeList) {
	var (
		s1, s2 *BUQNode
	)

	switch dir {
	case north:
		s1 = n.northEast.(*BUQNode)
		s2 = n.northWest.(*BUQNode)
		break
	case east:
		s1 = n.northEast.(*BUQNode)
		s2 = n.southEast.(*BUQNode)
		break
	case south:
		s1 = n.southEast.(*BUQNode)
		s2 = n.southWest.(*BUQNode)
		break
	case west:
		s1 = n.northWest.(*BUQNode)
		s2 = n.southWest.(*BUQNode)
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
// cf. Hanan Samet 1981 article Neighbor Finding in Quadtrees.
// It can return nil if the neighbor can't be found.
func (n *BUQNode) equalSizeNeighbor(dir side) *BUQNode {
	var neighbor *BUQNode

	// Ascent the tree up to a common ancestor.
	if n.parent != nil {
		buparent := n.parent.(*BUQNode)
		if adjacent(dir, n.location) {
			neighbor = buparent.equalSizeNeighbor(dir)
		} else {
			neighbor = buparent
		}
	}

	// Backtrack mirroring the ascending moves.
	if neighbor != nil && !neighbor.isLeaf() {
		return neighbor.child(reflect(dir, n.location)).(*BUQNode)
	}
	return neighbor
}

// neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *BUQNode) neighbours(dir side, nodes *QNodeList) {
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

// ForEachNeighbour calls the given function for each neighbour of current
// node.
func (n *BUQNode) ForEachNeighbour(fn func(QNode)) {
	var nodes QNodeList
	n.neighbours(north, &nodes)
	n.neighbours(south, &nodes)
	n.neighbours(east, &nodes)
	n.neighbours(west, &nodes)
	for _, nb := range nodes {
		fn(nb)
	}
}
