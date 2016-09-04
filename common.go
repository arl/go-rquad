package quadtree

//go:generate stringer -type=quadrant
type quadrant int

const (
	northWest quadrant = iota
	northEast
	southWest
	southEast
)

//go:generate stringer -type=side
type side int

const (
	north side = iota
	east
	south
	west
)

// Init initializes the package and should be called before using it.
func Init() {
	// initialize the quadrant-side adjacency array
	arrAdjacent = [4][4]bool{
		/*       NW     NE     SW     SE  */
		/* N */ {true, true, false, false},
		/* E */ {false, true, false, true},
		/* S */ {false, false, true, true},
		/* W */ {true, false, true, false},
	}

	// initialize the mirror-quadrant array
	arrReflect = [4][4]quadrant{
		/*           NW         NE         SW         SE    */
		/* N */ {southWest, southEast, northWest, northEast},
		/* E */ {northEast, northWest, southEast, southWest},
		/* S */ {southWest, southEast, northWest, northEast},
		/* W */ {northEast, northWest, southEast, southWest},
	}
}

var (
	arrAdjacent [4][4]bool
	arrReflect  [4][4]quadrant
)

// adjacent checks if a quadrant is adjacent to a given side of this node.
func adjacent(s side, q quadrant) bool {
	return arrAdjacent[s][q]
}

// reflect obtains the mirror image of a quadrant on a given side.
func reflect(s side, q quadrant) quadrant {
	return arrReflect[s][q]
}
