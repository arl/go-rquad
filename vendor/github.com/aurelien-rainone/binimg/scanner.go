package binimg

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

	// Uniform reports wether the given region is uniform. If is the case, the
	// uniform color bit is returned, otherwise the returned Bit is not
	// significative (always the zero value of Bit).
	Uniform(r image.Rectangle) (bool, color.Color)
}
