package imgscan

import (
	"bytes"
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

type binaryScanner struct {
	*binimg.Image
}

// IsUniformColor indicates if the region r is only made of pixels of color c.
//
// The scan stops at the first pixel encountered that is different from c.
func (s *binaryScanner) IsUniformColor(r image.Rectangle, c color.Color) bool {
	var (
		ok  bool
		bit binimg.Bit
	)
	// ensure c is a binimg.Bit, or convert it
	if bit, ok = c.(binimg.Bit); !ok {
		bit = s.ColorModel().Convert(c).(binimg.Bit)
	}

	// in a binary image, pixel/bytes are 1 or 0, we want the other color for
	// bytes.IndexBytes
	other := bit.Other().V
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

// IsUniform indicates if the region r is uniform. If that is the case, the
// uniform color is returned, otherwise the returned color is nil.
//
// The scan stops at the first pixel encountered that is different from the
// previous one.
func (s *binaryScanner) IsUniform(r image.Rectangle) (bool, color.Color) {
	// bit color of the first pixel (top-left)
	first := s.BitAt(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.IsUniformColor(r, first) {
		return true, first
	}
	return false, nil
}

// AverageColor indicates wether the region is uniform and the average color
// of the region r. If all the pixels have the same color (i.e the region is
// uniform) then the average color is that color.
//
// A full scan of the region is performed in order to determine the average
// color.
func (s *binaryScanner) AverageColor(r image.Rectangle) (bool, color.Color) {
	// if region is uniform
	if uniform, col := s.IsUniform(r); uniform {
		// return its color
		return true, col
	}
	// or consider the whole region as made of On pixel (arbitrary)
	return false, binimg.On
}

// NewBinaryScanner creates a binary scanner from a binary image.
func NewBinaryScanner(img *binimg.Image) Scanner {
	return &binaryScanner{img}
}
