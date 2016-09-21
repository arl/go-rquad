package quadtree

import (
	"fmt"
	"image"
)

// QNodeColor is the set of colors that can take a QNode.
type QNodeColor byte

const (
	// Black is the color of leaf nodes that
	// are considered as obstructed.
	Black QNodeColor = 0

	// White is the color of leaf nodes that
	// are considered as free.
	White = 1

	// Gray is the color of non-leaf nodes that
	// contain both black and white children.
	Gray = 2
)

// QNode defines the interface for a quadtree node.
type QNode interface {
	Parent() QNode

	NorthWest() QNode
	NorthEast() QNode
	SouthWest() QNode
	SouthEast() QNode

	TopLeft() image.Point
	BottomRight() image.Point

	Color() QNodeColor

	// Neighbours fills a NodeList with the neighbours of this node. n must be
	// a leaf node, or nodes will be an empty slice.
	Neighbours(nodes *NodeList)
}

// quadnode is a basic implementation of the QNode interface.
type quadnode struct {
	parent QNode // pointer to the parent node

	northWest QNode // pointer to the northwest child
	northEast QNode // pointer to the northeast child
	southWest QNode // pointer to the southwest child
	southEast QNode // pointer to the southeast child

	// node top-left corner coordinates, the origin
	topLeft image.Point

	// node bottom-right corner coordinates, the point is included
	bottomRight image.Point

	// node color
	color QNodeColor
}

func (n *quadnode) TopLeft() image.Point {
	return n.topLeft
}

func (n *quadnode) BottomRight() image.Point {
	return n.bottomRight
}

func (n *quadnode) Color() QNodeColor {
	return n.color
}

func (n *quadnode) NorthWest() QNode {
	return n.northWest
}

func (n *quadnode) NorthEast() QNode {
	return n.northEast
}

func (n *quadnode) SouthWest() QNode {
	return n.southWest
}

func (n *quadnode) SouthEast() QNode {
	return n.southEast
}

func (n *quadnode) Parent() QNode {
	return n.parent
}

func (n *quadnode) width() int {
	return n.bottomRight.X - n.topLeft.X
}

func (n *quadnode) height() int {
	return n.bottomRight.Y - n.topLeft.Y
}

// child returns a pointer to the child node associated to the given quadrant
func (n *quadnode) child(q quadrant) QNode {
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
