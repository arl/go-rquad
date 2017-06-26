package rquad

import (
	"errors"
	"fmt"
	"image"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

// Color is the set of colors that can take a Node.
type Color byte

const (
	// Black is the color of leaf nodes
	// that are considered as obstructed.
	Black Color = 0 + iota

	// White is the color of leaf nodes
	// that are considered as free.
	White

	// Gray is the color of non-leaf nodes
	// that contain both black and white children.
	Gray
)

const colorName = "BlackWhiteGray"

var colorIndex = [...]uint8{0, 5, 10, 14}

func (i Color) String() string {
	if i >= Color(len(colorIndex)-1) {
		return fmt.Sprintf("Color(%d)", i)
	}
	return colorName[colorIndex[i]:colorIndex[i+1]]
}

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

type ColoredNode struct {
	BasicNode
	color Color // node color
}

// Parent returns the quadtree node that is the parent of current one.
func (n *ColoredNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *ColoredNode) Child(q Quadrant) Node {
	if n.c[q] == nil {
		return nil
	}
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *ColoredNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *ColoredNode) Color() Color {
	return n.color
}

// IsLeaf returns wether the node is a leaf node
func (n *ColoredNode) IsLeaf() bool {
	return n.leaf
}

// Location returns the node inside its parent quadrant
func (n *ColoredNode) Location() Quadrant {
	return n.location
}

type NodeSetter interface {
	ScanAndSet(*Node)
	NewRoot() Node
	NewNode(Node, Quadrant, image.Rectangle) Node
}

type BinImageNodeSetter struct {
	scanner    imgscan.Scanner
	resolution int
}

// Create a NewBinImageSetter, that knows how to create quadtree nodes from a
// binary image
func NewBinImageSetter(scanner imgscan.Scanner, resolution int) (*BinImageNodeSetter, error) {
	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}

	// To ensure a consistent behavior and eliminate corner cases,
	// the Quadtree's root node needs to have children. Thus, the
	// first instantiated Node needs to always be subdivided.
	// This condition asserts the resolution is respected.
	minDim := scanner.Bounds().Dx()
	if scanner.Bounds().Dy() < minDim {
		minDim = scanner.Bounds().Dy()
	}
	if minDim < resolution*2 {
		return nil, errors.New("the image smaller dimension must be greater or equal to twice the resolution")
	}
	return &BinImageNodeSetter{
		scanner:    scanner,
		resolution: resolution,
	}, nil
}

func (s *BinImageNodeSetter) NewRoot() Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			leaf:   false,
			bounds: s.scanner.Bounds(),
		},
	}
}

func (s *BinImageNodeSetter) NewNode(parent Node, location Quadrant, bounds image.Rectangle) Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			bounds:   bounds,
			location: location,
			parent:   parent,
		},
	}
}

func (s *BinImageNodeSetter) ScanAndSet(n *Node) {
	colNode := (*n).(*ColoredNode)
	uniform, col := s.scanner.IsUniform((*colNode).bounds)
	switch uniform {
	case true:
		// quadrant is uniform, won't need to subdivide any further
		if col == binimg.White {
			colNode.color = White
		} else {
			colNode.color = Black
		}
		colNode.leaf = true
	case false:
		// if we reached maximal resolution..
		if (*colNode).bounds.Dx()/2 < s.resolution || (*colNode).bounds.Dy()/2 < s.resolution {
			// ...make this node a black leaf, instead of gray
			(*colNode).color = Black
			(*colNode).leaf = true
		} else {
			(*colNode).leaf = false
		}
	}
}
