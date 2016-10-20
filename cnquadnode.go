package quadtree

import "image"

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
	parent *CNQNode // pointer to the parent node
	// TODO: use an array for children
	northWest *CNQNode        // pointer to the northwest child
	northEast *CNQNode        // pointer to the northeast child
	southWest *CNQNode        // pointer to the southwest child
	southEast *CNQNode        // pointer to the southeast child
	bounds    image.Rectangle // node bounds
	color     Color           // node color
	cn        [4]*CNQNode     // cardinal neighbours
	location  quadrant        // node location inside its parent
	size      int             // size of a quadrant side
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *CNQNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color() returns the node Color.
func (n *CNQNode) Color() Color {
	return n.color
}

func (n *CNQNode) updateNECardinalNeighbours() {
	if n.parent == nil || n.cn[north] == nil {
		// nothing to update as this quadrant lies on the north border
		return
	}
	// step 2.2: Updating Cardinal Neighbors of NE sub-Quadrant.
	if n.cn[north] != nil {
		if n.cn[north].size < n.size {
			C0 := n.northWest
			C1 := n.northEast
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
			C0 := n.northWest
			C2 := n.southWest
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
		n.forEachNeighbour(west, func(qn Node) {
			western := qn.(*CNQNode)
			if western.cn[east] == n {
				if western.bounds.Max.Y > n.southWest.bounds.Min.Y {
					// choose SW
					western.cn[east] = n.southWest
				} else {
					// choose NW
					western.cn[east] = n.northWest
				}
				if western.cn[east].bounds.Min.Y == western.bounds.Min.Y {
					western.cn[east].cn[west] = western
				}
			}
		})
	}

	if n.cn[north] != nil {
		n.forEachNeighbour(north, func(qn Node) {
			northern := qn.(*CNQNode)
			if northern.cn[south] == n {
				if northern.bounds.Max.X > n.northEast.bounds.Min.X {
					// choose NE
					northern.cn[south] = n.northEast
				} else {
					// choose NW
					northern.cn[south] = n.northWest
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
			n.cn[east].cn[west] = n.northEast
		}
	}

	if n.cn[south] != nil {
		// To update the southern CN of a quadrant Q that is being
		// decomposed: Q.CN3.CN1=Q.Ch[SE]
		// TODO: this seems a typo in the paper.
		// should have read this instead: Q.CN3.CN1=Q.Ch[SW]
		if n.cn[south] != nil && n.cn[south].cn[north] == n {
			n.cn[south].cn[north] = n.southWest
		}
	}
}

// Location() returns the node inside its parent quadrant
func (n *CNQNode) Location() quadrant {
	return n.location
}

// Parent returns the quadtree node that is the parent of current one.
func (n *CNQNode) Parent() Node {
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *CNQNode) Child(q quadrant) Node {
	switch q {
	case northWest:
		return n.northWest
	case northEast:
		return n.northEast
	case southWest:
		return n.southWest
	default:
		fallthrough
	case southEast:
		return n.southEast
	}
}

// forEachNeighbour calls fn on every neighbour of the current node in the given
// direction
func (n *CNQNode) forEachNeighbour(dir side, fn func(Node)) {
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
func (n *CNQNode) ForEachNeighbour(fn func(Node)) {
	n.forEachNeighbour(west, fn)
	n.forEachNeighbour(north, fn)
	n.forEachNeighbour(east, fn)
	n.forEachNeighbour(south, fn)
}
