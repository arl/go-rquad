package rquad

import "image"

// PointLocator is the interface implemented by objects that, given a point in
// 2D space, can return the leaf node that contains it.
type PointLocator interface {
	// PointLocation returns the Node that contains the given point, or nil.
	PointLocation(image.Point) Node
}

// PointLocation returns the quadtree node containing the given point.
//
// The standard method to search for the Node that contains a given point is a
// recursively search, starting from the root node, returning the Node that
// contains the point is a leaf. If q implements the PointLocator interfacen
// then the specific and more efficient implementation of PointLocation is
// called
func PointLocation(q Quadtree, pt image.Point) Node {
	if locator, ok := q.(PointLocator); ok {
		// use the specific point location implementation
		return locator.PointLocation(pt)
	}
	return pointLocation(q.Root(), pt)
}

func pointLocation(n Node, pt image.Point) Node {
	if !pt.In(n.Bounds()) {
		return nil
	}
	if n.Color() != Gray {
		return n
	}

	if pt.In(n.Child(Northwest).Bounds()) {
		return pointLocation(n.Child(Northwest), pt)
	} else if pt.In(n.Child(Northeast).Bounds()) {
		return pointLocation(n.Child(Northeast), pt)
	} else if pt.In(n.Child(Southwest).Bounds()) {
		return pointLocation(n.Child(Southwest), pt)
	}
	return pointLocation(n.Child(Southeast), pt)
}
