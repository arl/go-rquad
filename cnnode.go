package rquad

import "image"

// CNNode is a node of a CNTree.
//
// It is an implementation of the Node interface, with additional fields and
// methods required to obtain the node neighbours in constant time. The time
// complexity reduction is obtained through the addition of only four pointers per
// leaf node in the quadtree.
//
// - The Western cardinal neighbor is the top-most neighbor node among the
//   western neighbors, noted cn0.
// - The Northern cardinal neighbor is the left-most neighbor node among the
//   northern neighbors, noted cn1.
// - The Eastern cardinal neighbor is the bottom-most neighbor node among the
//   eastern neighbors, noted cn2.
// - The Southern cardinal neighbor is the right-most neighbor node among the
//   southern neighbors, noted cn3.
type CNNode struct {
	parent   *CNNode         // pointer to the parent node
	c        [4]*CNNode      // children nodes
	cn       [4]*CNNode      // cardinal neighbours
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
	size     int             // size of a quadrant side
	id       int             // unique identifier
}

// Parent returns the quadtree node that is the parent of current one.
func (n *CNNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *CNNode) Child(q Quadrant) Node {
	if n.c[q] == nil {
		return nil
	}
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *CNNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *CNNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *CNNode) Location() Quadrant {
	return n.location
}

// Id returns an identifier for this node, guaranteed to be unique inside a
// Quadtree.
func (n *CNNode) Id() int {
	return n.id
}

func (n *CNNode) updateNorthEast() {
	if n.parent == nil || n.cn[North] == nil {
		// nothing to update as this quadrant lies on the north border
		return
	}
	// step 2.2: Updating Cardinal Neighbors of NE sub-Quadrant.
	if n.cn[North] != nil {
		if n.cn[North].size < n.size {
			c0 := n.c[Northwest]
			c0.cn[North] = n.cn[North]
			// to update C1, we perform a west-east traversal
			// recording the cumulative size of traversed nodes
			cur := c0.cn[North]
			cumsize := cur.size
			for cumsize < c0.size {
				cur := cur.cn[East]
				cumsize += cur.size
			}
			n.c[Northeast].cn[North] = cur
		}
	}
}

func (n *CNNode) updateSouthWest() {
	if n.parent == nil || n.cn[West] == nil {
		// nothing to update as this quadrant lies on the west border
		return
	}
	// step 2.1: Updating Cardinal Neighbors of SW sub-Quadrant.
	if n.cn[North] != nil {
		if n.cn[North].size < n.size {
			c0 := n.c[Northwest]
			c0.cn[North] = n.cn[North]
			// to update C2, we perform a north-south traversal
			// recording the cumulative size of traversed nodes
			cur := c0.cn[West]
			cumsize := cur.size
			for cumsize < c0.size {
				cur := cur.cn[South]
				cumsize += cur.size
			}
			n.c[Southwest].cn[West] = cur
		}
	}
}

// updateNeighbours updates all neighbours according to the current
// decomposition.
func (n *CNNode) updateNeighbours() {
	// On each direction, a full traversal of the neighbors
	// should be performed.  In every quadrant where a reference
	// to the parent quadrant is stored as the Cardinal Neighbor,
	// it should be replaced by one of its children created after
	// the decomposition

	if n.cn[West] != nil {
		n.forEachNeighbour(West, func(qn Node) {
			western := qn.(*CNNode)
			if western.cn[East] == n {
				if western.bounds.Max.Y > n.c[Southwest].bounds.Min.Y {
					// choose SW
					western.cn[East] = n.c[Southwest]
				} else {
					// choose NW
					western.cn[East] = n.c[Northwest]
				}
				if western.cn[East].bounds.Min.Y == western.bounds.Min.Y {
					western.cn[East].cn[West] = western
				}
			}
		})
	}

	if n.cn[North] != nil {
		n.forEachNeighbour(North, func(qn Node) {
			northern := qn.(*CNNode)
			if northern.cn[South] == n {
				if northern.bounds.Max.X > n.c[Northeast].bounds.Min.X {
					// choose NE
					northern.cn[South] = n.c[Northeast]
				} else {
					// choose NW
					northern.cn[South] = n.c[Northwest]
				}
				if northern.cn[South].bounds.Min.X == northern.bounds.Min.X {
					northern.cn[South].cn[North] = northern
				}
			}
		})
	}

	if n.cn[East] != nil {
		if n.cn[East] != nil && n.cn[East].cn[West] == n {
			// To update the eastern CN of a quadrant Q that is being
			// decomposed: Q.CN2.CN0=Q.Ch[NE]
			n.cn[East].cn[West] = n.c[Northeast]
		}
	}

	if n.cn[South] != nil {
		// To update the southern CN of a quadrant Q that is being
		// decomposed: Q.CN3.CN1=Q.Ch[SE]
		// TODO: there seems to be a typo in the paper.
		// should have read this instead: Q.CN3.CN1=Q.Ch[SW]
		if n.cn[South] != nil && n.cn[South].cn[North] == n {
			n.cn[South].cn[North] = n.c[Southwest]
		}
	}
}

// forEachNeighbour calls fn on every neighbour of the current node in the given
// direction.
func (n *CNNode) forEachNeighbour(dir Side, fn func(Node)) {
	// start from the cardinal neighbour on the given direction
	N := n.cn[dir]
	if N == nil {
		return
	}
	fn(N)
	if N.size >= n.size {
		return
	}

	traversal := traversal(dir)
	opposite := opposite(dir)
	// perform cardinal neighbour traversal
	for {
		N = N.cn[traversal]
		if N != nil && N.cn[opposite] == n {
			fn(N)
		} else {
			return
		}
	}
}

// ForEachNeighbour calls the given function for each neighbour of current
// node.
func (n *CNNode) ForEachNeighbour(fn func(Node)) {
	n.forEachNeighbour(West, fn)
	n.forEachNeighbour(North, fn)
	n.forEachNeighbour(East, fn)
	n.forEachNeighbour(South, fn)
}
