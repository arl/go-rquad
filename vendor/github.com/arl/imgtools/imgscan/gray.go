package imgscan

import (
	"bytes"
	"image"
	"image/color"
)

type grayScanner struct {
	*image.Gray
}

// IsUniformColor indicates if the region r is only made of pixels of color c.
//
// The scan stops at the first pixel encountered that is different from c.
func (s *grayScanner) IsUniformColor(r image.Rectangle, c color.Color) bool {
	var (
		ok   bool       // conversion to color.Gray ok
		gray color.Gray // c converted to Gray
		b    uint8      // gray level
		last int        // index of last byte to check on the line
	)
	// ensure c is a color.Gray, or convert it
	if gray, ok = c.(color.Gray); !ok {
		gray = s.ColorModel().Convert(c).(color.Gray)
	}
	b, last = gray.Y, r.Max.X-r.Min.X-1

	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.PixOffset(r.Min.X, y)
		j := s.PixOffset(r.Max.X, y)
		// check the first and the last pixel/byte of color c are respectively
		// the first and last pixels of the slice.
		if bytes.IndexByte(s.Pix[i:j], b) != 0 || bytes.LastIndexByte(s.Pix[i:j], b) != last {
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
func (s *grayScanner) IsUniform(r image.Rectangle) (bool, color.Color) {
	// gray color of the first pixel (top-left)
	first := s.GrayAt(r.Min.X, r.Min.Y)

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
func (s *grayScanner) AverageColor(r image.Rectangle) (bool, color.Color) {
	if uniform, col := s.IsUniform(r); uniform {
		return true, col
	}

	var sum uint64
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			sum += uint64(s.GrayAt(x, y).Y)
		}
	}
	return false, color.Gray{uint8(sum / uint64(r.Dx()*r.Dy()))}
}

// NewGrayScanner creates a gray scanner from a gray image.
func NewGrayScanner(img *image.Gray) Scanner {
	return &grayScanner{img}
}
