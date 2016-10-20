package quadtree

import "image"

// Color is the set of colors that can take a QNode.
type Color byte

const (
	// Black is the color of leaf nodes
	// that are considered as obstructed.
	Black Color = 0 + iota

	// White is the color of leaf nodes
	// that are considered as free.
	White

	// Gray is the color of non-leaf nodes
	// that contain both black and white children.
	Gray
)

// QNode defines the interface for a quadtree node.
type QNode interface {

	// Parent returns the quadtree node that is the parent of current one.
	Parent() QNode

	// Child returns current node child at specified quadrant.
	Child(quadrant) QNode

	// Bounds returns the bounds of the rectangular area represented by this
	// quadtree node.
	Bounds() image.Rectangle

	// Color() returns the node Color.
	Color() Color

	// Location() returns the node inside its parent quadrant
	Location() quadrant
}

// AdjacencyNode is a QNode that can find its adjacent nodes, or neighbours.
type AdjacencyNode interface {
	QNode
	// ForEachNeighbour calls the given function
	// for each neighbour of current node.
	// TODO: rename ForEachAdjacent
	ForEachNeighbour(func(QNode))
}
