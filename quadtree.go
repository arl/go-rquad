package quadtree

import "image"

// QNodeList is a slice of QNode's
type QNodeList []QNode

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// ForEachLeaf calls the given function for each leaf node of the quadtree.
	//
	// Successive calls to the provided function are performed in no particular
	// order. The color parameter allows to loop on the leaves of a particular
	// color, Black or White.
	// NOTE: As by definition, Gray leaves do not exist, passing Gray to
	// ForEachLeaf should return all leaves, independently of their color.
	ForEachLeaf(Color, func(QNode))

	// Root returns the quadtree root node.
	Root() QNode
}

// PointLocator is the interface implemented by objects having a PointLocation method.
type PointLocator interface {
	// PointLocation returns the quadtree node containing the given point.
	PointLocation(image.Point) QNode
}

// CodeLocator is the interface implemented by objects having a CodeLocation method.
type CodeLocator interface {
	// CodeLocation returns the quadtree node corresponding to a given location code.
	CodeLocation(uint64) QNode
}
