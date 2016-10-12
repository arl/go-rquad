package quadtree

import "fmt"

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
	qnode             // base quadnode
	cn0      *CNQNode // western cardinal neighbour
	cn1      *CNQNode // northern cardinal neighbour
	cn2      *CNQNode // eastern cardinal neighbour
	cn3      *CNQNode // southern cardinal neighbour
	location quadrant // node location
	size     int      // size of quadrant sides
}

func (n *CNQNode) updateNECardinalNeighbours() {
	if n.parent == nil || n.cn1 == nil {
		// nothing to update as this quadrant lies on the north border
		return
	}
	// step 2.2: Updating Cardinal Neighbors of NE sub-Quadrant.
	C0 := n.northWest.(*CNQNode)
	C1 := n.northEast.(*CNQNode)

	if n.cn1 != nil {
		if n.cn1.size >= n.size {
			C0.cn1 = n.cn1
			C1.cn1 = n.cn1
		} else {
			C0.cn1 = n.cn1
			// to update C1, we perform a west-east traversal
			cur := C0.cn1
			// TODO: here we could initialize cumsize with cur.size and avoid
			// to enter in the loop if not needed
			cumsize := 0 // cumulative size of traversed cardinal neighbours
			for cumsize < C0.size {
				cumsize += cur.size
				tmp := cur.cardinalNeighbour(east)
				if tmp == nil {
					break
				}
				cur = tmp
			}
			C1.cn1 = cur
		}
	}
}

func (n *CNQNode) updateSWCardinalNeighbours() {
	if n.parent == nil || n.cn0 == nil {
		// nothing to update as this quadrant lies on the west border
		return
	}
	// step 2.1: Updating Cardinal Neighbors of SW sub-Quadrant.
	C0 := n.northWest.(*CNQNode)
	C2 := n.southWest.(*CNQNode)
	if n.cn1 != nil {
		if n.cn1.size >= n.size {
			C0.cn0 = n.cn0
			C2.cn0 = n.cn0
		} else {
			C0.cn1 = n.cn1
			// to update C1, we perform a north-south traversal
			cur := C0.cn0
			// TODO: here we could initialize cumsize with cur.size and avoid
			// to enter in the loop if not needed
			cumsize := 0 // cumulative size of traversed cardinal neighbours
			for cumsize < C0.size {
				cumsize += cur.size
				tmp := cur.cardinalNeighbour(south)
				if tmp == nil {
					break
				}
				cur = tmp
			}
			C2.cn0 = cur
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
		if western.cn2 == n {
			if western.bounds.Max.Y > SW.bounds.Min.Y {
				// choose SW
				western.cn2 = SW
			} else {
				// choose NW
				western.cn2 = NW
			}
			if western.cn2.bounds.Min.Y == western.bounds.Min.Y {
				western.cn2.cn0 = western
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
		if northern.cn3 == n {
			if northern.bounds.Max.X > NE.bounds.Min.X {
				// choose NE
				northern.cn3 = NE
			} else {
				// choose NW
				northern.cn3 = NW
			}
			if northern.cn3.bounds.Min.X == northern.bounds.Min.X {
				northern.cn3.cn1 = northern
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

	if n.cn2 != nil && n.cn2.cn0 == n {
		// parent is stored as the cn
		n.cn2.cn0 = n.northEast.(*CNQNode)
	}
}

// Step3UpdateSouth updates the southern neighbours of current quadrant.
func (n *CNQNode) Step3UpdateSouth() {
	// To update the southern CN of a quadrant Q that is being
	// decomposed: Q.CN3.CN1=Q.Ch[SE]
	// TODO: could the paper be wrong about that?
	// and mean this instead: Q.CN3.CN1=Q.Ch[SW]
	if n.cn3 != nil && n.cn3.cn1 == n {

		//Step3UpdateSouth of  [(2,2)-2-Gray|CN ←(0,2)-2 ↑(2,0)-2 →(4,0)-4 ↓(0,4)-4]
		//n.cn3 was [(0,4)-4-White|CN ← ↑(0,2)-2 →(4,4)-4 ↓]
		//n.cn3  is [(0,4)-4-White|CN ← ↑(2,3)-1 →(4,4)-4 ↓]
		n.cn3.cn1 = n.southWest.(*CNQNode)
	}
}

// isLeaf checks if this node is a leaf, i.e. is either black or white.
func (n *CNQNode) isLeaf() bool {
	return n.color != Gray
}

func (n *CNQNode) cardinalNeighbour(dir side) *CNQNode {
	// TODO: should use an array for cardinal neighbours so we can index them
	//       so we won't need this function but just to do n.cn[0]
	switch dir {
	case west:
		return n.cn0
	case north:
		return n.cn1
	case east:
		return n.cn2
	default:
		fallthrough
	case south:
		return n.cn3
	}
}

// neighbours locates all leaf neighbours of the current node in the given
// direction, appending them to a slice.
func (n *CNQNode) neighbours(dir side, nodes *QNodeList) {
	switch dir {

	case north:
		N := n.cardinalNeighbour(north)
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform west to east traversal
				for {
					N = N.cardinalNeighbour(east)
					if N != nil && N.cn3 == n {
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
		N := n.cardinalNeighbour(west)
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform north to south traversal
				for {
					N = N.cardinalNeighbour(south)
					if N != nil && N.cn2 == n {
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
		N := n.cardinalNeighbour(south)
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform east to west traversal
				for {
					N = N.cardinalNeighbour(west)
					if N != nil && N.cn1 == n {
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
		N := n.cardinalNeighbour(east)
		if N != nil {
			*nodes = append(*nodes, N)
			if N.size < n.size {
				// perform south to north traversal
				for {
					N = N.cardinalNeighbour(north)
					if N != nil && N.cn0 == n {
						*nodes = append(*nodes, N)
					} else {
						break
					}
				}

			}
		}
	}
}

// Neighbours returns the node neighbours. n should be
// a leaf node, or the returned slice will be empty.
func (n *CNQNode) Neighbours(nodes *QNodeList) {
	var _n, _s, _e, _w QNodeList
	n.neighbours(north, &_n)
	n.neighbours(south, &_s)
	n.neighbours(east, &_e)
	n.neighbours(west, &_w)
	*nodes = append(*nodes, _n...)
	*nodes = append(*nodes, _s...)
	*nodes = append(*nodes, _e...)
	*nodes = append(*nodes, _w...)
}

func (n *CNQNode) String() string {
	var scn0, scn1, scn2, scn3 string
	if n.cn0 != nil {
		scn0 = fmt.Sprintf("%v-%d", n.cn0.bounds.Min, n.cn0.size)
	}
	if n.cn1 != nil {
		scn1 = fmt.Sprintf("%v-%d", n.cn1.bounds.Min, n.cn1.size)
	}
	if n.cn2 != nil {
		scn2 = fmt.Sprintf("%v-%d", n.cn2.bounds.Min, n.cn2.size)
	}
	if n.cn3 != nil {
		scn3 = fmt.Sprintf("%v-%d", n.cn3.bounds.Min, n.cn3.size)
	}
	return fmt.Sprintf("[%v-%d-%s|CN ←%v ↑%v →%v ↓%v]", n.bounds.Min, n.size, n.color, scn0, scn1, scn2, scn3)
}
