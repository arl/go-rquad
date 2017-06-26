package rquad

import (
	"image"
)

// Node defines the interface for a quadtree node.
type Node interface {

	// SetChild set the child node at specified quadrant.
	SetChild(q Quadrant, n Node)

	// Bounds returns the bounds of the rectangular area represented by this
	// quadtree node.
	Bounds() image.Rectangle

	IsLeaf() bool
}

// NodeList is a slice of Node instances.
type NodeList []Node
