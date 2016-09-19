package quadtree

import "image"

//go:generate stringer -type=NodeColor
type NodeColor byte

const (
	Black NodeColor = 0
	White           = 1
	Gray            = 2
)

// Quadnode defines the interface for a quadtree node.
type Quadnode interface {
	Parent() Quadnode

	NorthWest() Quadnode
	NorthEast() Quadnode
	SouthWest() Quadnode
	SouthEast() Quadnode

	TopLeft() image.Point
	BottomRight() image.Point

	Color() NodeColor
}
