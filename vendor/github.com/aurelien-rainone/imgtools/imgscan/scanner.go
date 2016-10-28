package imgscan

import (
	"image"
	"image/color"
)

// A Scanner is an image that can report wether a rectangular region is uniform
// (i.e composed of the same color) or not.
type Scanner interface {
	image.Image

	// UniformColor reports wether all the pixels of given region are of the color c.
	UniformColor(r image.Rectangle, c color.Color) bool

	// Uniform reports wether the region specified by r is uniform or not, and if
	// that is the case the uniform color is returned. If the region is not
	// uniform, the color is undefined.
	Uniform(r image.Rectangle) (uniform bool, col color.Color)
}
