package binimg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

// NewScanner returns a new Scanner of the given image.Image.
//
// The actual scanner implementation depends on the image bit depth and the
// availability of an implementation.
func NewScanner(img image.Image) (Scanner, error) {
	switch impl := img.(type) {
	case *Binary:
		return &binaryScanner{impl}, nil
	case *image.Alpha:
	case *image.Gray:
	default:
	}
	// NOTE:
	// an efficient scanner for images using a 8bit depth color model would
	// be easy to write if there was an efficient function in the Go
	// standard bytes package that was similar to the C++
	// std::find_first_not_of function (i.e returns the index of the first
	// byte of a slice that different from a given byte, or a set of bytes)
	return nil, fmt.Errorf("unsupported image type")
}

type binaryScanner struct {
	*Binary
}

// UniformColor reports wether all the pixels of given region are of the color c.
func (s *binaryScanner) UniformColor(r image.Rectangle, c color.Color) bool {
	// we want the other color for bytes.IndexBytes
	other := c.(Bit).Other().v
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.PixOffset(r.Min.X, y)
		j := s.PixOffset(r.Max.X, y)
		// look for the first pixel that is not c
		if bytes.IndexByte(s.Pix[i:j], other) != -1 {
			return false
		}
	}
	return true
}

// Uniform reports wether the given region is uniform. If is the case, the
// uniform color bit is returned, otherwise the returned Bit is not
// significative (always the zero value of Bit).
func (s *binaryScanner) Uniform(r image.Rectangle) (bool, color.Color) {
	// bit color of the first pixel (top-left)
	first := s.BitAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, Bit{}
}
