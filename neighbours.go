package rquad

// NeighbourNode is a Node that can access to its neighbours, by neighbour it
// is intended other leaf nodes of the same color.
type NeighbourNode interface {
	Node
	// ForEachNeighbour calls the given function
	// for each neighbour of current node.
	ForEachNeighbour(func(Node))
}

// ForEachNeighbour calls the given function for each neighbour of the quadtree
// node n.
//
// The neighbour finding technique used depends on the Node implementation. If
// the node implements the NeighbourNode interface, then the specific and
// faster implementation of ForEachNeighbour is called. If that's not the case,
// the neighbours are found by using the generic but slower "bottom-up
// neighbour finding technique", cf. Hanan Samet 1981 article Neighbour Finding
// in Quadtrees
func ForEachNeighbour(n Node, fn func(Node)) {
	if adjnode, ok := n.(NeighbourNode); ok {
		// use neighbour node specific implementation
		adjnode.ForEachNeighbour(fn)
		return
	}

	// perform generic implementation (bottom-up technique)
	neighbours(n, North, fn)
	neighbours(n, South, fn)
	neighbours(n, East, fn)
	neighbours(n, West, fn)
}

// equalSizeNeighbour locates an equal-sized neighbour of the current node in the
// vertical or horizontal direction.
//
// cf. Hanan Samet 1981 article Neighbour Finding in Quadtrees.
// It can return nil if the neighbour can't be found.
func equalSizeNeighbour(n Node, dir Side) Node {
	var neighbour Node

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

// neighbours calls fn for each leaf neighbours of the current node it finds in
// the given direction
func neighbours(n Node, dir Side, fn func(Node)) {
	// If no neighbour can be found in the given
	// direction, node will be null.
	node := equalSizeNeighbour(n, dir)
	if node != nil {
		if node.Color() != Gray {
			// Neighbour is already a leaf node, we're done after that.
			fn(node)
		} else {
			// The neighbour isn't a leaf node so we need to
			// go further down matching its children, but in
			// the opposite direction from where we came.
			children(node, opposite(dir), fn)
		}
	}
}

// children calls fn for each leaf children of this node it finds in the given
// direction.
func children(n Node, dir Side, fn func(Node)) {
	var (
		s1, s2 Node
	)

	switch dir {
	case North:
		s1 = n.Child(Northeast)
		s2 = n.Child(Northwest)
		break
	case East:
		s1 = n.Child(Northeast)
		s2 = n.Child(Southeast)
		break
	case South:
		s1 = n.Child(Southeast)
		s2 = n.Child(Southwest)
		break
	case West:
		s1 = n.Child(Northwest)
		s2 = n.Child(Southwest)
	}

	if s1.Color() != Gray {
		fn(s1)
	} else {
		children(s1, dir, fn)
	}

	if s2.Color() != Gray {
		fn(s2)
	} else {
		children(s2, dir, fn)
	}
}
