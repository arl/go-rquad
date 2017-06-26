package rquad

import (
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

type BinImgTreeModel struct {
	*ImageTreeModel
}

func NewBinImgTreeModel(scanner imgscan.Scanner, resolution int) (*BinImgTreeModel, error) {
	model, err := NewImageTreeModel(scanner, resolution)
	return &BinImgTreeModel{model}, err
}

func (m *BinImgTreeModel) ScanAndSet(n *Node) {
	colNode := (*n).(*ColoredNode)
	uniform, col := m.scanner.IsUniform((*colNode).bounds)
	switch uniform {
	case true:
		// quadrant is uniform, won't need to subdivide any further
		if col == binimg.White {
			colNode.color = color.White
		} else {
			colNode.color = color.Black
		}
		colNode.leaf = true
	case false:
		// if we reached maximal resolution..
		if colNode.bounds.Dx()/2 < m.resolution || colNode.bounds.Dy()/2 < m.resolution {
			// ...make this node a black leaf, instead of gray
			colNode.color = color.Black
			colNode.leaf = true
		} else {
			colNode.leaf = false
		}
	}
}
