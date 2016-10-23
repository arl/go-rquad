// Package quadtrees proposes various implementations of region quadtrees.
// The region quadtree is a special kind of quadtree that recursively
// subdivides a 2D dimensional space into 4, smaller and generally equal
// rectangular regions, until the wanted quadtree resolution has been reached,
// or no further subdivisions can be performed.
//
// Region quadtree may be used for image processing, in this case a node
// represents a rectangular region of an image in which all pixels have the
// same color.
//
// A region quadtree may also be used as a variable resolution representation
// of a data field. For example, the temperatures in an area may be stored as a
// quadtree, with each leaf node storing the average temperature over the
// subregion it represents.
//
// Quadtree implementations in this package use the binimg.Scanner interface to
// represent the complete area, and provide us with a way to know if a
// particular sub-area is to be considered uniform, in which case further
// subdivision is not necessary, or not.

package rquad

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// ForEachLeaf calls the given function for each leaf node of the quadtree.
	//
	// Successive calls to the provided function are performed in no particular
	// order. The color parameter allows to loop on the leaves of a particular
	// color, Black or White.
	// NOTE: As by definition, Gray leaves do not exist, passing Gray to
	// ForEachLeaf should return all leaves, independently of their color.
	ForEachLeaf(Color, func(Node))

	// Root returns the quadtree root node.
	Root() Node
}
