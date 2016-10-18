package quadtree

import "image"

// QNodeColor is the set of colors that can take a QNode.
type QNodeColor byte

const (
	// Black is the color of leaf nodes
	// that are considered as obstructed.
	Black QNodeColor = 0 + iota

	// White is the color of leaf nodes
	// that are considered as free.
	White

	// Gray is the color of non-leaf nodes
	// that contain both black and white children.
	Gray
)

// QNode defines the interface for a quadtree node.
type QNode interface {
	Parent() QNode

	NorthWest() QNode
	NorthEast() QNode
	SouthWest() QNode
	SouthEast() QNode

	Bounds() image.Rectangle
	Color() QNodeColor

	// ForEachNeighbour calls the given function for each neighbour of current
	// node.
	ForEachNeighbour(func(QNode))
}

// qnode is a basic implementation of the QNode interface and serves as
// an embeddable struct to various QNode implementations.
type qnode struct {
	parent QNode // pointer to the parent node

	northWest QNode // pointer to the northwest child
	northEast QNode // pointer to the northeast child
	southWest QNode // pointer to the southwest child
	southEast QNode // pointer to the southeast child

	// node bounds
	bounds image.Rectangle

	// node color
	color QNodeColor

	// node location for its parent
	location quadrant
}

func (n *qnode) Color() QNodeColor {
	return n.color
}

func (n *qnode) NorthWest() QNode {
	return n.northWest
}

func (n *qnode) NorthEast() QNode {
	return n.northEast
}

func (n *qnode) SouthWest() QNode {
	return n.southWest
}

func (n *qnode) SouthEast() QNode {
	return n.southEast
}

func (n *qnode) Parent() QNode {
	return n.parent
}

func (n *qnode) Bounds() image.Rectangle {
	return n.bounds
}

// child returns a pointer to the child node associated to the given quadrant
func (n *qnode) child(q quadrant) QNode {
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
