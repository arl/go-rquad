package quadtree

import (
	"image"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
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

	Color() bmp.Color
}
