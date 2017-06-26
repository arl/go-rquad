package rquad

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

type BinImgTreeModel struct {
	scanner    imgscan.Scanner
	resolution int
}

func NewBinImgTreeModel(scanner imgscan.Scanner, resolution int) (*BinImgTreeModel, error) {
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
	return &BinImgTreeModel{
		scanner:    scanner,
		resolution: resolution,
	}, nil
}

func (m *BinImgTreeModel) NewRoot() Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			leaf:   false,
			bounds: m.scanner.Bounds(),
		},
	}
}

func (m *BinImgTreeModel) NewNode(parent Node, location Quadrant, bounds image.Rectangle) Node {
	return &ColoredNode{
		BasicNode: BasicNode{
			bounds:   bounds,
			location: location,
			parent:   parent,
		},
	}
}

func (m *BinImgTreeModel) ScanAndSet(n *Node) {
	colNode := (*n).(*ColoredNode)
	uniform, col := m.scanner.IsUniform((*colNode).bounds)
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
		if colNode.bounds.Dx()/2 < m.resolution || colNode.bounds.Dy()/2 < m.resolution {
			// ...make this node a black leaf, instead of gray
			colNode.color = Black
			colNode.leaf = true
		} else {
			colNode.leaf = false
		}
	}
}
