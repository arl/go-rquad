package quadtree

// Quadrant indicates the position of a child Node inside its parent.
type Quadrant int

const (
	Northwest Quadrant = iota
	Northeast
	Southwest
	Southeast
	rootQuadrant
)

type Side int

const (
	West Side = iota
	North
	East
	South
)

// init() initializes package level variables.
func init() {

	// initialize the quadrant-side adjacency array
	arrAdjacent = [4][4]bool{
		/*       NW     NE     SW     SE  */
		/* W */ {true, false, true, false},
		/* N */ {true, true, false, false},
		/* E */ {false, true, false, true},
		/* S */ {false, false, true, true},
	}

	// initialize the mirror-quadrant array
	arrReflect = [4][4]Quadrant{
		/*           NW         NE         SW         SE    */
		/* W */ {Northeast, Northwest, Southeast, Southwest},
		/* N */ {Southwest, Southeast, Northwest, Northeast},
		/* E */ {Northeast, Northwest, Southeast, Southwest},
		/* S */ {Southwest, Southeast, Northwest, Northeast},
	}

	// initialize the opposite sides array
	arrOpposite = [4]Side{
		/* W     N      E     S  */
		East, South, West, North,
	}

	// For Cardinal Neighbour Quadtrees
	arrTraversal = [4]Side{
		/* W     N      E     S  */
		South, East, North, West,
	}
}

var (
	arrAdjacent  [4][4]bool
	arrReflect   [4][4]Quadrant
	arrOpposite  [4]Side
	arrTraversal [4]Side
)

// adjacent checks if a quadrant is adjacent to a given side of this node.
func adjacent(s Side, q Quadrant) bool {
	return arrAdjacent[s][q]
}

// reflect obtains the mirror image of a quadrant on a given side.
func reflect(s Side, q Quadrant) Quadrant {
	return arrReflect[s][q]
}

// opposite returns the opposite of a side.
func opposite(s Side) Side {
	return arrOpposite[s]
}

// traversal returns for a given cardinal neighbour direction,
// the direction of the neighbour traversal.
func traversal(s Side) Side {
	return arrTraversal[s]
}
