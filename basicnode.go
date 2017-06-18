package rquad

import "image"

// BasicNode represents a standard quadtree node.
//
// It is a basic implementation of the Node interface, the one used in the
// BasicTree implementation of the Quadtree interface.
type BasicNode struct {
	parent   Node            // pointer to the parent node
	c        [4]Node         // children nodes
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
}

// Parent returns the quadtree node that is the parent of current one.
func (n *BasicNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *BasicNode) Child(q Quadrant) Node {
	if n.c[q] == nil {
		return nil
	}
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *BasicNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *BasicNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *BasicNode) Location() Quadrant {
	return n.location
}
