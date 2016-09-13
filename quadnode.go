package quadtree

import (
	"image"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

// quadnode is a basic implementation of the Quadnode interface.
type quadnode struct {
	parent Quadnode // pointer to the parent node

	northWest Quadnode // pointer to the northwest child
	northEast Quadnode // pointer to the northeast child
	southWest Quadnode // pointer to the southwest child
	southEast Quadnode // pointer to the southeast child

	// node top-left corner coordinates, the origin
	topLeft image.Point

	// node bottom-right corner coordinates, the point is included
	bottomRight image.Point

	// node color
	color bmp.Color
}

func (n *quadnode) TopLeft() image.Point {
	return n.topLeft
}

func (n *quadnode) BottomRight() image.Point {
	return n.bottomRight
}

func (n *quadnode) Color() bmp.Color {
	return n.color
}

func (n *quadnode) NorthWest() Quadnode {
	return n.northWest
}

func (n *quadnode) NorthEast() Quadnode {
	return n.northEast
}

func (n *quadnode) SouthWest() Quadnode {
	return n.southWest
}

func (n *quadnode) SouthEast() Quadnode {
	return n.southEast
}

func (n *quadnode) Parent() Quadnode {
	return n.parent
}

func (n *quadnode) width() int {
	return n.bottomRight.X - n.topLeft.X
}

func (n *quadnode) height() int {
	return n.bottomRight.Y - n.topLeft.Y
}

// child returns a pointer to the child node associated to the given quadrant
func (n *quadnode) child(q quadrant) Quadnode {
	switch q {
	case northWest:
		return n.northWest
	case southWest:
		return n.southWest
	case northEast:
		return n.northEast
	case southEast:
		return n.southEast
	}
	panic("undefined quadrant")
}

// quadrant obtains this node's quadrant relative to its parent.
//
// must not be called on the root node
func (n *quadnode) quadrant() quadrant {
	if n.parent == nil {
		panic("the root node's quadrant is undefined")
	}

	if n.parent.NorthWest() == n {
		return northWest
	} else if n.parent.SouthWest() == n {
		return southWest
	} else if n.parent.NorthEast() == n {
		return northEast
	} else {
		return southEast
	}
}

// inbound checks if a given point is inside the region represented by this
// node.
func (n *quadnode) inbound(pt image.Point) bool {
	return (n.topLeft.X <= pt.X && pt.X <= n.bottomRight.X) &&
		(n.topLeft.Y <= pt.Y && pt.Y <= n.bottomRight.Y)
}
