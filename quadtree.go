package quadtree

import "image"

// QNodeList is a slice of QNode's
type QNodeList []QNode

// Quadtree defines the interface for a quadtree type.
type Quadtree interface {
	// WhiteNodes returns a slice of all the white nodes of the quadtree.
	WhiteNodes() QNodeList

	// Root returns the quadtree root node.
	Root() QNode
}

// Query returns the leaf node that contains a given point.
func Query(q Quadtree, pt image.Point) (n QNode, exists bool) {
	return query(q.Root(), pt)
}

func query(n QNode, pt image.Point) (QNode, bool) {
	if !pt.In(n.Bounds()) {
		return nil, false
	}
	if n.Color() != Gray {
		return n, true
	}

	if pt.In(n.NorthWest().Bounds()) {
		return query(n.NorthWest(), pt)
	} else if pt.In(n.NorthEast().Bounds()) {
		return query(n.NorthEast(), pt)
	} else if pt.In(n.SouthWest().Bounds()) {
		return query(n.SouthWest(), pt)
	}
	return query(n.SouthEast(), pt)
}
