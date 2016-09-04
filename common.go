package quadtree

// Init initializes the package and should be called before using it.
func Init() {
	// initialize the adjacent sides array
	adjacentSides = [4][4]bool{
		/*       NW     NE     SW     SE  */
		/* N */ {true, true, false, false},
		/* E */ {false, true, false, true},
		/* S */ {false, false, true, true},
		/* W */ {true, false, true, false},
	}
}

var adjacentSides [4][4]bool

// adjacent checks if a quadrant is adjacent to a given side of this node.
func adjacent(s side, q quadrant) bool {
	return adjacentSides[s][q]
}

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
