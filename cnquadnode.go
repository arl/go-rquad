package quadtree

// CNQNode is a node of a CNQuadtree.
//
// It is an implementation of the QNode interface, with additional fields and
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
type CNQNode struct {
	qnode                // base quadnode
	cn       [4]*CNQNode // cardinal neighbours
	location quadrant    // node location
	size     int         // size of quadrant sides
}

func (n *CNQNode) updateNECardinalNeighbours() {
	if n.parent == nil || n.cn[north] == nil {
		// nothing to update as this quadrant lies on the north border
		return
	}
	// step 2.2: Updating Cardinal Neighbors of NE sub-Quadrant.
	if n.cn[north] != nil {
		if n.cn[north].size < n.size {
			C0 := n.northWest.(*CNQNode)
			C1 := n.northEast.(*CNQNode)
			C0.cn[north] = n.cn[north]
			// to update C1, we perform a west-east traversal
			// recording the cumulative size of traversed nodes
			cur := C0.cn[north]
			cumsize := cur.size
			for cumsize < C0.size {
				tmp := cur.cn[east]
				if tmp == nil {
					break
				}
				cur = tmp
				cumsize += cur.size
			}
			C1.cn[north] = cur
		}
	}
}

func (n *CNQNode) updateSWCardinalNeighbours() {
	if n.parent == nil || n.cn[west] == nil {
		// nothing to update as this quadrant lies on the west border
		return
	}
	// step 2.1: Updating Cardinal Neighbors of SW sub-Quadrant.
	if n.cn[north] != nil {
		if n.cn[north].size < n.size {
			C0 := n.northWest.(*CNQNode)
			C2 := n.southWest.(*CNQNode)
			C0.cn[north] = n.cn[north]
			// to update C2, we perform a north-south traversal
			// recording the cumulative size of traversed nodes
			cur := C0.cn[west]
			cumsize := cur.size
			for cumsize < C0.size {
				tmp := cur.cn[south]
				if tmp == nil {
					break
				}
				cur = tmp
				cumsize += cur.size
			}
			C2.cn[west] = cur
		}
	}
}

// updateNeighbours updates all neighbours according to the current
// decomposition.
func (n *CNQNode) updateNeighbours() {
	if n.cn[west] != nil {
		n.forEachNeighbour(west, func(qn QNode) {
			western := qn.(*CNQNode)
			if western.cn[east] == n {
				if western.bounds.Max.Y > n.southWest.(*CNQNode).bounds.Min.Y {
					// choose SW
					western.cn[east] = n.southWest.(*CNQNode)
				} else {
					// choose NW
					western.cn[east] = n.northWest.(*CNQNode)
				}
				if western.cn[east].bounds.Min.Y == western.bounds.Min.Y {
					western.cn[east].cn[west] = western
				}
			}
		})
	}

	if n.cn[north] != nil {
		n.forEachNeighbour(north, func(qn QNode) {
			northern := qn.(*CNQNode)
			if northern.cn[south] == n {
				if northern.bounds.Max.X > n.northEast.(*CNQNode).bounds.Min.X {
					// choose NE
					northern.cn[south] = n.northEast.(*CNQNode)
				} else {
					// choose NW
					northern.cn[south] = n.northWest.(*CNQNode)
				}
				if northern.cn[south].bounds.Min.X == northern.bounds.Min.X {
					northern.cn[south].cn[north] = northern
				}
			}
		})
	}

	if n.cn[east] != nil {
		if n.cn[east] != nil && n.cn[east].cn[west] == n {
			// To update the eastern CN of a quadrant Q that is being
			// decomposed: Q.CN2.CN0=Q.Ch[NE]
			n.cn[east].cn[west] = n.northEast.(*CNQNode)
		}
	}

	if n.cn[south] != nil {
		// To update the southern CN of a quadrant Q that is being
		// decomposed: Q.CN3.CN1=Q.Ch[SE]
		// TODO: this seems a typo in the paper.
		// should have read this instead: Q.CN3.CN1=Q.Ch[SW]
		if n.cn[south] != nil && n.cn[south].cn[north] == n {
			n.cn[south].cn[north] = n.southWest.(*CNQNode)
		}
	}
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *CNQNode) isLeaf() bool {
	return n.color != Gray
}

// forEachNeighbour calls fn on every neighbour of the current node in the given
// direction
func (n *CNQNode) forEachNeighbour(dir side, fn func(QNode)) {
	// start from the cardinal neighbour on the given direction
	N := n.cn[dir]
	if N == nil {
		return
	}
	fn(N)
	if N.size >= n.size {
		return
	}

	switch dir {

	case west:
		// perform north to south traversal
		for {
			N = N.cn[south]
			if N != nil && N.cn[east] == n {
				fn(N)
			} else {
				return
			}
		}

	case north:
		// perform west to east traversal
		for {
			N = N.cn[east]
			if N != nil && N.cn[south] == n {
				fn(N)
			} else {
				break
			}
		}

	case east:
		// perform south to north traversal
		for {
			N = N.cn[north]
			if N != nil && N.cn[west] == n {
				fn(N)
			} else {
				return
			}
		}

	case south:
		// perform east to west traversal
		for {
			N = N.cn[west]
			if N != nil && N.cn[north] == n {
				fn(N)
			} else {
				return
			}
		}
	}
}

// ForEachNeighbour calls the given function for each neighbour of current
// node.
func (n *CNQNode) ForEachNeighbour(fn func(QNode)) {
	n.forEachNeighbour(west, fn)
	n.forEachNeighbour(north, fn)
	n.forEachNeighbour(east, fn)
	n.forEachNeighbour(south, fn)
}
