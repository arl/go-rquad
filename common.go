package rquad

import "fmt"

// Quadrant indicates the position of a child Node inside its parent.
type Quadrant int

// Possible values for the Quadrant type.
const (
	Northwest Quadrant = iota
	Northeast
	Southwest
	Southeast
	rootQuadrant
)

const _Quadrant_name = "NorthwestNortheastSouthwestSoutheastrootQuadrant"

var _Quadrant_index = [...]uint8{0, 9, 18, 27, 36, 48}

func (i Quadrant) String() string {
	if i < 0 || i >= Quadrant(len(_Quadrant_index)-1) {
		return fmt.Sprintf("Quadrant(%d)", i)
	}
	return _Quadrant_name[_Quadrant_index[i]:_Quadrant_index[i+1]]
}

// Side is used to represent a direction according to a quadtree Node.
type Side int

// Possible values for the Side type.
const (
	West Side = iota
	North
	East
	South
)

const _Side_name = "WestNorthEastSouth"

var _Side_index = [...]uint8{0, 4, 9, 13, 18}

func (i Side) String() string {
	if i < 0 || i >= Side(len(_Side_index)-1) {
		return fmt.Sprintf("Side(%d)", i)
	}
	return _Side_name[_Side_index[i]:_Side_index[i+1]]
}

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

	// for Cardinal Neighbour Quadtrees
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
