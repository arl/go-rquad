package quadtree

import "image"

// BUQNode is a node of a BUQuadtree.
//
// It is a basic implementation of the Node interface, augmented with
// methods implementing the bottom-up neighbour finding techniques.
type BUQNode struct {
	parent   *BUQNode        // pointer to the parent node
	c        [4]*BUQNode     // children nodes
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *BUQNode) Bounds() image.Rectangle {
	return n.bounds
}

// Parent returns the quadtree node that is the parent of current one.
func (n *BUQNode) Parent() Node {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent
}

// Color returns the node Color.
func (n *BUQNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *BUQNode) Location() Quadrant {
	return n.location
}

// Child returns current node child at specified quadrant.
func (n *BUQNode) Child(q Quadrant) Node {
	return n.c[q]
}
