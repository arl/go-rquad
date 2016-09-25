package binimg

import (
	"image"
	"image/color"
	"image/draw"
)

// pow2roundup rounds up to next higher power
// of 2, or n if n is already a power of 2.
func pow2roundup(x int) int {
	if x <= 1 {
		return 1
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

// PowerOf2Image returns a power-of-2 square image, on which src is copied over
// at the origin point {0,0}.
//
// The new image has the dimensions of the smallest square that can contain the
// whole src image. The process creates	an image of uniform color, with pad, and
// copies the original image. pad is converted to src color.Model so that the
// returned image has the same color model as the original.
// Note: if src dimensions is already a power-of-2 square image, it is returned
// as-is.
func PowerOf2Image(src image.Image, pad color.Color) image.Image {
	if IsPowerOf2Image(src) {
		return src
	}

	var (
		dst  *image.RGBA
		side int // square side
	)

	side = src.Bounds().Dx()
	if src.Bounds().Dy() > side {
		side = src.Bounds().Dy()
	}
	side = pow2roundup(side)

	// compute the dimensions
	x, y := src.Bounds().Min.X, src.Bounds().Min.Y

	// create a uniform square image at those dimensions
	dst = image.NewRGBA(image.Rect(x, y, x+side, y+side))
	cpad := src.ColorModel().Convert(pad)
	draw.Draw(dst, dst.Bounds(), &image.Uniform{cpad}, image.ZP, draw.Src)

	// now draw the original image onto it
	draw.Draw(dst, src.Bounds(), src, image.ZP, draw.Src)
	return dst
}

// IsPowerOf2Image reports wether img is a power-of-2 square image or not.
func IsPowerOf2Image(img image.Image) bool {
	maxdim := img.Bounds().Dx()
	if img.Bounds().Dy() > maxdim {
		maxdim = img.Bounds().Dy()
	}
	maxdim = pow2roundup(maxdim)
	return maxdim == img.Bounds().Dx() &&
		maxdim == img.Bounds().Dy()
}
