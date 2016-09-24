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

// Query allows to obtain the maximum-depth node that contains a point.
func Query(q Quadtree, pt image.Point) (n QNode, exists bool) {
	return query(q.Root(), pt)
}

func query(n QNode, pt image.Point) (QNode, bool) {
	if !Inbound(n, pt) {
		return nil, false
	}
	if n.Color() != Gray {
		return n, true
	}

	if Inbound(n.NorthWest(), pt) {
		return query(n.NorthWest(), pt)
	} else if Inbound(n.NorthEast(), pt) {
		return query(n.NorthEast(), pt)
	} else if Inbound(n.SouthWest(), pt) {
		return query(n.SouthWest(), pt)
	}
	return query(n.SouthEast(), pt)
}

// Inbound checks if a given point is inside
// the region represented by this node.
func Inbound(n QNode, pt image.Point) bool {
	return (n.TopLeft().X <= pt.X && pt.X < n.BottomRight().X) &&
		(n.TopLeft().Y <= pt.Y && pt.Y < n.BottomRight().Y)
}
