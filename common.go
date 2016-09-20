package quadtree

type quadrant int

const (
	northWest quadrant = iota
	northEast
	southWest
	southEast
)

type side int

const (
	north side = iota
	east
	south
	west
)

var initDone bool

// initPackage initializes package level variables.
func initPackage() {
	if initDone {
		return
	}

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

	// initialize the opposite sides array
	arrOpposite = [4]side{
		/* N     E      S     W  */
		south, west, north, east,
	}

	initDone = true
}

var (
	arrAdjacent [4][4]bool
	arrReflect  [4][4]quadrant
	arrOpposite [4]side
)

// adjacent checks if a quadrant is adjacent to a given side of this node.
func adjacent(s side, q quadrant) bool {
	return arrAdjacent[s][q]
}

// reflect obtains the mirror image of a quadrant on a given side.
func reflect(s side, q quadrant) quadrant {
	return arrReflect[s][q]
}

// opposite returns, given a side, its opposite
func opposite(s side) side {
	return arrOpposite[s]
}
