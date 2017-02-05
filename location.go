package rquad

import "image"

// PointLocator is the interface implemented by objects that, given a point in
// 2D space, can return the leaf node that contains it.
type PointLocator interface {
	// PointLocation returns the Node that contains the given point, or nil.
	PointLocation(image.Point) Node
}

// PointLocation returns the leaf node containing the given point.
//
// The generic method to search for the leaf Node that contains a given point is
// a recursive search from the root node, it returns the leaf node containing
// the point.
// If q implements the PointLocator interface, (i.e it implements a probably
// more efficient method), then the call is forwardded to q.PointLocation
func PointLocation(q Quadtree, pt image.Point) Node {
	if locator, ok := q.(PointLocator); ok {
		// use the specific point location implementation
		return locator.PointLocation(pt)
	}
	return pointLocation(q.Root(), pt)
}

// generic recursive method to return the leaf node containing pt
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
