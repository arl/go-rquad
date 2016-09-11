package quadtree

import (
	"image"

	"github.com/RookieGameDevs/quadtree/bmp"
)

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

type Subdivider interface {
	Subdivide(Quadnode)
}

type QuadtreeCreator interface {
	CreateRootNode() Quadnode
	CreateInnerNode(topleft, bottomright image.Point, parent Quadnode) Quadnode
	Subdivider
}
