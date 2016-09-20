package quadtree

import "image"

// NodeList is a slice of Quadnode's
type NodeList []Quadnode

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// PointQuery returns the Quadnode containing the point at given coordinates.
	//
	// If such node doesn't exist, exists is false
	PointQuery(pt image.Point) (n Quadnode, exists bool)
}
