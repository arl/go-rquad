package quadtree

import (
	"fmt"
	"image"
)

type NodeColor byte

const (
	Black NodeColor = 0
	White           = 1
	Gray            = 2
)

// Quadnode defines the interface for a quadtree node.
type Quadnode interface {
	Parent() Quadnode

	NorthWest() Quadnode
	NorthEast() Quadnode
	SouthWest() Quadnode
	SouthEast() Quadnode

	TopLeft() image.Point
	BottomRight() image.Point

	Color() NodeColor

	// Neighbours fills a NodeList with the neighbours of this node. n must be
	// a leaf node, or nodes will be an empty slice.
	Neighbours(nodes *NodeList)
}

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
	color NodeColor
}

func (n *quadnode) TopLeft() image.Point {
	return n.topLeft
}

func (n *quadnode) BottomRight() image.Point {
	return n.bottomRight
}

func (n *quadnode) Color() NodeColor {
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
	default:
		return n.southEast
	}
}

// inbound checks if a given point is inside the region represented by this
// node.
func (n *quadnode) inbound(pt image.Point) bool {
	return (n.topLeft.X <= pt.X && pt.X < n.bottomRight.X) &&
		(n.topLeft.Y <= pt.Y && pt.Y < n.bottomRight.Y)
}

func (n *quadnode) String() string {
	return fmt.Sprintf("(%d,%d %d,%d %s)", n.topLeft.X, n.topLeft.Y, n.bottomRight.X, n.bottomRight.Y, n.color)
}
