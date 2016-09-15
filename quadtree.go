package quadtree

import "image"

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// PointQuery returns the Quadnode containing the point at given coordinates.
	//
	// If such node doesn't exist, exists is false
	PointQuery(pt image.Point) (n Quadnode, exists bool)
}
