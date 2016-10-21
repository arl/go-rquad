package rquad

import "image"

// CNode is a node of a CNTree.
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
//
// - 𝜌(𝐷) returns the immediate parent of the node D. The notation 𝜌²(𝐷)
//   denotes the parent of the parent of D. 𝜌°(𝐷) = 𝐷.
// - 𝑆𝑖𝑧𝑒(𝐷) returns the side length of node N in pixels.
// - 𝜑 𝑖(𝐷) returns the cardinal Neighbor of node D in direction i,
//   for 𝑖 ∈  0,1,2,3 where 0,1,2,3 represent respectively the directions West,
//   North, East and South.
// - 𝜑 𝑖𝑗(𝐷) represents the Cardinal Neighbor in the direction i of the
//   Cardinal Neighbor in direction j of the Node D. 𝜑 𝑖𝑗(𝐷) = 𝜑 𝑖(𝜑 𝑗(𝐷))
// - 𝜑 𝑖(𝜑 𝑖(𝐷)) will be noted as 𝜑 𝑖²(𝐷). This represents the Cardinal
//   Neighbor in the direction i of the Cardinal Neighbor in direction i of the
//   Node D for 𝑖 ∈ 0,1,2,3 where 0,1,2,3 represent respectively the directions
//   West, North, East and South and where 𝜑 𝑖°(𝐷)=𝐷. 𝜑 𝑖²(𝐷) = 𝜑 𝑖(𝜑 𝑖(𝐷 ))
type cnNode struct {
	parent   *cnNode         // pointer to the parent node
	c        [4]*cnNode      // children nodes
	cn       [4]*cnNode      // cardinal neighbours
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
	size     int             // size of a quadrant side
}

// Parent returns the quadtree node that is the parent of current one.
func (n *cnNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *cnNode) Child(q Quadrant) Node {
	if n.c[q] == nil {
		return nil
	}
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *cnNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *cnNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *cnNode) Location() Quadrant {
	return n.location
}

func (n *cnNode) updateNECardinalNeighbours() {
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
				tmp := cur.cn[East]
				if tmp == nil {
					break
				}
				cur = tmp
				cumsize += cur.size
			}
			n.c[Northeast].cn[North] = cur
		}
	}
}

func (n *cnNode) updateSWCardinalNeighbours() {
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
				tmp := cur.cn[South]
				if tmp == nil {
					break
				}
				cur = tmp
				cumsize += cur.size
			}
			n.c[Southwest].cn[West] = cur
		}
	}
}

// updateNeighbours updates all neighbours according to the current
// decomposition.
func (n *cnNode) updateNeighbours() {
	if n.cn[West] != nil {
		n.forEachNeighbour(West, func(qn Node) {
			western := qn.(*cnNode)
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
			northern := qn.(*cnNode)
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
		// TODO: this seems a typo in the paper.
		// should have read this instead: Q.CN3.CN1=Q.Ch[SW]
		if n.cn[South] != nil && n.cn[South].cn[North] == n {
			n.cn[South].cn[North] = n.c[Southwest]
		}
	}
}

// forEachNeighbour calls fn on every neighbour of the current node in the given
// direction
func (n *cnNode) forEachNeighbour(dir Side, fn func(Node)) {
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
func (n *cnNode) ForEachNeighbour(fn func(Node)) {
	n.forEachNeighbour(West, fn)
	n.forEachNeighbour(North, fn)
	n.forEachNeighbour(East, fn)
	n.forEachNeighbour(South, fn)
}
