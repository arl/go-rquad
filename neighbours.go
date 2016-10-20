package quadtree

// equalSizeNeighbour locates an equal-sized neighbour of the current node in the
// vertical or horizontal direction.
//
// cf. Hanan Samet 1981 article Neighbour Finding in Quadtrees.
// It can return nil if the neighbour can't be found.
func equalSizeNeighbour(n QNode, dir side) QNode {
	var neighbour QNode

	// Ascent the tree up to a common ancestor.
	parent := n.Parent()
	if parent != nil {
		if adjacent(dir, n.Location()) {
			neighbour = equalSizeNeighbour(parent, dir)
		} else {
			neighbour = parent
		}
	}

	// Backtrack mirroring the ascending moves.
	if neighbour != nil && neighbour.Color() == Gray {
		return neighbour.Child(reflect(dir, n.Location()))
	}
	return neighbour
}

// neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func neighbours(n QNode, dir side, nodes *QNodeList) {
	// If no neighbour can be found in the given
	// direction, node will be null.
	node := equalSizeNeighbour(n, dir)
	if node != nil {
		if node.Color() != Gray {
			// Neighbour is already a leaf node, we're done.
			*nodes = append(*nodes, node)
		} else {
			// The neighbour isn't a leaf node so we need to
			// go further down matching its children, but in
			// the opposite direction from where we came.
			children(node, opposite(dir), nodes)
		}
	}
}

// children fills the given slice with all the leaf children of this node (i.e
// either black or white), that can be found in a given direction.
func children(n QNode, dir side, nodes *QNodeList) {
	var (
		s1, s2 QNode
	)

	switch dir {
	case north:
		s1 = n.Child(northEast)
		s2 = n.Child(northWest)
		break
	case east:
		s1 = n.Child(northEast)
		s2 = n.Child(southEast)
		break
	case south:
		s1 = n.Child(southEast)
		s2 = n.Child(southWest)
		break
	case west:
		s1 = n.Child(northWest)
		s2 = n.Child(southWest)
	}

	if s1.Color() != Gray {
		*nodes = append(*nodes, s1)
	} else {
		children(s1, dir, nodes)
	}

	if s2.Color() != Gray {
		*nodes = append(*nodes, s2)
	} else {
		children(s2, dir, nodes)
	}
}

// ForEachNeighbour calls the given function for each neighbour of the quadtree
// node n.
//
// The neighbour finding technique used depends on the QNode implementation. If
// the node implements the AdjacencyNode interface, then the specific and
// faster implementation of ForEachNeighbour is called. If that's not the case,
// the neighbours are found by using the generic but slower "bottom-up
// neighbour finding technique", cf. Hanan Samet 1981 article Neighbour Finding
// in Quadtrees
func ForEachNeighbour(n QNode, fn func(QNode)) {
	if adjnode, ok := n.(AdjacencyNode); ok {
		// use adjacency node specific implementation
		adjnode.ForEachNeighbour(fn)
		return
	}

	// TODO; fn should be passed to individual neighbours functions to remove
	// the need to fill a temporary slice.
	var nodes QNodeList
	neighbours(n, north, &nodes)
	neighbours(n, south, &nodes)
	neighbours(n, east, &nodes)
	neighbours(n, west, &nodes)
	for _, nb := range nodes {
		fn(nb)
	}
}
