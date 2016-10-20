package quadtree

import "image"

// BUQNode is a node of a BasicTree.
//
// It is a basic implementation of the Node interface, augmented with
// methods implementing the bottom-up neighbour finding techniques.
type BasicNode struct {
	parent   *BasicNode      // pointer to the parent node
	c        [4]*BasicNode   // children nodes
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *BasicNode) Bounds() image.Rectangle {
	return n.bounds
}

// Parent returns the quadtree node that is the parent of current one.
func (n *BasicNode) Parent() Node {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent
}

// Color returns the node Color.
func (n *BasicNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *BasicNode) Location() Quadrant {
	return n.location
}

// Child returns current node child at specified quadrant.
func (n *BasicNode) Child(q Quadrant) Node {
	return n.c[q]
}
