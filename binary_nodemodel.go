package rquad

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

type BinaryNodeModel struct {
	scanner    imgscan.Scanner
	resolution int
}

// Create a NewBinImageSetter, that knows how to create quadtree nodes from a
// binary image
func NewBinaryNodeModel(scanner imgscan.Scanner, resolution int) (*BinaryNodeModel, error) {
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
	return &BinaryNodeModel{
		scanner:    scanner,
		resolution: resolution,
	}, nil
}

func (s *BinaryNodeModel) NewRoot() Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			leaf:   false,
			bounds: s.scanner.Bounds(),
		},
	}
}

func (s *BinaryNodeModel) NewNode(parent Node, location Quadrant, bounds image.Rectangle) Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			bounds:   bounds,
			location: location,
			parent:   parent,
		},
	}
}

func (s *BinaryNodeModel) ScanAndSet(n *Node) {
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
		if colNode.bounds.Dx()/2 < s.resolution || colNode.bounds.Dy()/2 < s.resolution {
			// ...make this node a black leaf, instead of gray
			colNode.color = Black
			colNode.leaf = true
		} else {
			colNode.leaf = false
		}
	}
}
