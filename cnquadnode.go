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
	C0 := n.northWest.(*CNQNode)
	C1 := n.northEast.(*CNQNode)

	if n.cn[north] != nil {
		if n.cn[north].size >= n.size {
			C0.cn[north] = n.cn[north]
			C1.cn[north] = n.cn[north]
		} else {
			C0.cn[north] = n.cn[north]
			// to update C1, we perform a west-east traversal
			cur := C0.cn[north]
			// TODO: here we could initialize cumsize with cur.size and avoid
			// to enter in the loop if not needed
			cumsize := 0 // cumulative size of traversed cardinal neighbours
			for cumsize < C0.size {
				cumsize += cur.size
				tmp := cur.cn[east]
				if tmp == nil {
					break
				}
				cur = tmp
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
	C0 := n.northWest.(*CNQNode)
	C2 := n.southWest.(*CNQNode)
	if n.cn[north] != nil {
		if n.cn[north].size >= n.size {
			C0.cn[west] = n.cn[west]
			C2.cn[west] = n.cn[west]
		} else {
			C0.cn[north] = n.cn[north]
			// to update C1, we perform a north-south traversal
			cur := C0.cn[west]
			// TODO: here we could initialize cumsize with cur.size and avoid
			// to enter in the loop if not needed
			cumsize := 0 // cumulative size of traversed cardinal neighbours
			for cumsize < C0.size {
				cumsize += cur.size
				tmp := cur.cn[south]
				if tmp == nil {
					break
				}
				cur = tmp
			}
			C2.cn[west] = cur
		}
	}
}

// Step3UpdateWest updates the western neighbours of current quadrant.
func (n *CNQNode) Step3UpdateWest() {
	NW := n.northWest.(*CNQNode)
	SW := n.southWest.(*CNQNode)

	// TODO: change for a direct loop on the western neighbours
	var westernNeighbours QNodeList
	n.neighbours(west, &westernNeighbours)
	for _, neighbour := range westernNeighbours {
		western := neighbour.(*CNQNode)
		if western.cn[east] == n {
			if western.bounds.Max.Y > SW.bounds.Min.Y {
				// choose SW
				western.cn[east] = SW
			} else {
				// choose NW
				western.cn[east] = NW
			}
			if western.cn[east].bounds.Min.Y == western.bounds.Min.Y {
				western.cn[east].cn[west] = western
			}
		}
	}
}

// Step3UpdateNorth updates the northern neighbours of current quadrant.
func (n *CNQNode) Step3UpdateNorth() {
	NW := n.northWest.(*CNQNode)
	NE := n.northEast.(*CNQNode)

	// TODO: change for a direct loop on the northern neighbours
	var northernNeighbours QNodeList
	n.neighbours(north, &northernNeighbours)
	for _, neighbour := range northernNeighbours {
		northern := neighbour.(*CNQNode)
		if northern.cn[south] == n {
			if northern.bounds.Max.X > NE.bounds.Min.X {
				// choose NE
				northern.cn[south] = NE
			} else {
				// choose NW
				northern.cn[south] = NW
			}
			if northern.cn[south].bounds.Min.X == northern.bounds.Min.X {
				northern.cn[south].cn[north] = northern
			}
		}
	}
}

// Step3UpdateEast updates the eastern neighbours of current quadrant.
func (n *CNQNode) Step3UpdateEast() {
	// To update the eastern CN of a quadrant Q that is being
	// decomposed: Q.CN2.CN0=Q.Ch[NE]

	// On each direction, a full traversal of the neighbors should
	//be performed. In every quadrant where a reference to the
	//parent quadrant is stored as the Cardinal Neighbor, it
	//should be replaced by one of its children created after the
	//decomposition.To minimize the effort, the step 3 and step
	//2 will be performed in a single traversal on each side.

	if n.cn[east] != nil && n.cn[east].cn[west] == n {
		// parent is stored as the cn
		n.cn[east].cn[west] = n.northEast.(*CNQNode)
	}
}

// Step3UpdateSouth updates the southern neighbours of current quadrant.
func (n *CNQNode) Step3UpdateSouth() {
	// To update the southern CN of a quadrant Q that is being
	// decomposed: Q.CN3.CN1=Q.Ch[SE]
	// TODO: could the paper be wrong about that?
	// and mean this instead: Q.CN3.CN1=Q.Ch[SW]
	if n.cn[south] != nil && n.cn[south].cn[north] == n {
		n.cn[south].cn[north] = n.southWest.(*CNQNode)
	}
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *CNQNode) isLeaf() bool {
	return n.color != Gray
}

// neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *CNQNode) neighbours(dir side, nodes *QNodeList) {
	switch dir {

	case north:
		N := n.cn[north]
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform west to east traversal
				for {
					N = N.cn[east]
					if N != nil && N.cn[south] == n {
						*nodes = append(*nodes, N)
					} else {
						break
					}
				}
			}
		}

	case west:
		// On the western side, the neighbors are found starting
		// from the western CN and moving to the south.
		N := n.cn[west]
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform north to south traversal
				for {
					N = N.cn[south]
					if N != nil && N.cn[east] == n {
						*nodes = append(*nodes, N)
					} else {
						break
					}
				}
			}
		}

	case south:
		// for the southern side, the neighbors are identified
		// starting from the southern CN and moving to the west
		N := n.cn[south]
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform east to west traversal
				for {
					N = N.cn[west]
					if N != nil && N.cn[north] == n {
						*nodes = append(*nodes, N)
					} else {
						break
					}
				}
			}
		}

	case east:
		// For the eastern side, the neighbors are identified
		// starting from the Eastern CN and moving north
		N := n.cn[east]
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform south to north traversal
				for {
					N = N.cn[north]
					if N != nil && N.cn[west] == n {
						*nodes = append(*nodes, N)
					} else {
						break
					}
				}

			}
		}
	}
}

// ForEachNeighbour calls the given function for each neighbour of current
// node.
func (n *CNQNode) ForEachNeighbour(fn func(QNode)) {
	var nodes QNodeList
	n.neighbours(north, &nodes)
	n.neighbours(south, &nodes)
	n.neighbours(east, &nodes)
	n.neighbours(west, &nodes)
	for _, nb := range nodes {
		fn(nb)
	}
}
