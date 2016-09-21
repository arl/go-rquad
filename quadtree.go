package quadtree

import "image"

// QNodeList is a slice of QNode's
type QNodeList []QNode

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// PointQuery returns the QNode containing the point at given coordinates.
	//
	// If such node doesn't exist, exists is false
	PointQuery(pt image.Point) (n QNode, exists bool)
}
