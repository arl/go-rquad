package imgscan

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"github.com/aurelien-rainone/imgtools/binimg"
)

// NewScanner returns a new Scanner of the given image.Image.
//
// The actual scanner implementation depends on the image bit depth and the
// availability of an implementation.
func NewScanner(img image.Image) (Scanner, error) {
	switch impl := img.(type) {
	case *binimg.Binary:
		return &binaryScanner{impl}, nil
	case *image.Alpha:
	case *image.Gray:
	default:
	}
	return nil, fmt.Errorf("unsupported image type")
}

type binaryScanner struct {
	*binimg.Binary
}

// UniformColor reports wether all the pixels of given region are of the color c.
func (s *binaryScanner) UniformColor(r image.Rectangle, c color.Color) bool {
	// we want the other color for bytes.IndexBytes
	// Bit zero value is Off
	var other binimg.Bit
	if s.Palette.OffColor == c {
		other = binimg.On
	}

	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := s.PixOffset(r.Min.X, y)
		j := s.PixOffset(r.Max.X, y)
		if bytes.IndexByte(s.Pix[i:j], byte(other)) != -1 {
			// quit at the first byte that is not 'other'
			return false
		}
	}
	return true
}

// Uniform reports wether the region specified by r is uniform or not, and if
// that is the case the uniform color is returned. If the region is not
// uniform, the color is undefined.
func (s *binaryScanner) Uniform(r image.Rectangle) (bool, color.Color) {
	// color of the first pixel (top-left)
	first := s.At(r.Min.X, r.Min.Y)

	// check if all the pixels of the region are of this color.
	if s.UniformColor(r, first) {
		return true, first
	}
	return false, nil
}
