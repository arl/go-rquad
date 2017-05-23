package imgscan

import (
	"bytes"
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

type grayScanner struct {
	*image.Gray
}

// UniformColor reports wether all the pixels of given region are of the color c.
func (s *grayScanner) UniformColor(r image.Rectangle, c color.Color) bool {
	cb := c.(color.Gray).Y
	last := r.Max.Y - r.Min.Y - 1
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.PixOffset(r.Min.X, y)
		j := s.PixOffset(r.Max.X, y)
		// check the first and the last pixel/byte of color c are respectively
		// the first and last pixel of the slice.
		if bytes.IndexByte(s.Pix[i:j], cb) != 0 || bytes.LastIndexByte(s.Pix[i:j], cb) != last {
			return false
		}
	}
	return true
}

// IsUniform reports wether the given region is uniform. If that is the case, the
// uniform color is returned, otherwise the returned color should be ignored
// as it is not significative.
func (s *grayScanner) IsUniform(r image.Rectangle) (bool, color.Color) {
	// gray color of the first pixel (top-left)
	first := s.GrayAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, binimg.Bit{}
}
