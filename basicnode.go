package rquad

import "image"

// basicNode represents a standard quadtree node.
//
// It is a basic implementation of the Node interface, the one used in the
// BasicTree implementation of the Quadtree interface.
type basicNode struct {
	parent   *basicNode      // pointer to the parent node
	c        [4]*basicNode   // children nodes
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
}

// Parent returns the quadtree node that is the parent of current one.
func (n *basicNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *basicNode) Child(q Quadrant) Node {
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *basicNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *basicNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *basicNode) Location() Quadrant {
	return n.location
}
