// Package binimg proposes an in-memory binary image format, that is an image
// that has only two possible values for each pixel. Typically, the two colors
// used for a binary image are black and white, though any two colors can be
// used. Such images are also referred to as "bi-level", or "two-level".
//
// Binary implements the standard Go image.Image and draw.Image interfaces
// and embeds a two colors palette of type Palette, that itself implements the
// color.Model interface. Palette allows any color.Color to be converted to
// OffColor or OnColor.
//
// A pixel could be stored as a single bit, but as the main goal of this
// package is fast manipulation of binary images, Bit, the underlying pixel
// data type manipulated by Binary image, is 1 byte wide.
//
// Binary are instantiated by the following functions:
//
//  func New(r image.Rectangle, p Palette) *Binary
//  func NewFromImage(src image.Image, p Palette) *Binary
//
// Author: Aurélien Rainone
//
package binimg

import (
	"image"
	"image/color"
	"image/draw"
)

// Bit represents a 1-bit binary color.
type Bit uint8

// On and Off are the only two values that can take a Bit.
//go:generate stringer -type=Bit
var (
	Off = Bit(0)
	On  = Bit(255)
)

// RGBA returns the red, green, blue and alpha values for a Bit color.
//
// alpha is always 0xffff (fully opaque) and r, g, b are all 0 or all 0xffff.
// Note: a Bit is not mean to be directly converted to RGBA with this method,
// but through the binary Palette of a Binary image.
func (bit Bit) RGBA() (r, g, b, a uint32) {
	v := uint32(bit)
	v |= v << 8
	return v, v, v, 0xffff
}

// Binary is an in-memory image whose At method returns Bit values.
type Binary struct {
	// Pix holds the image's pixels, as 0 or 1 uint8 values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Palette is the image binary Palette
	Palette Palette
}

// ColorModel returns the image.Image's color model.
func (b *Binary) ColorModel() color.Model { return b.Palette }

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (b *Binary) Bounds() image.Rectangle { return b.Rect }

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (b *Binary) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(b.Rect)) {
		return b.Palette.OffColor
	}
	if b.BitAt(x, y) == Off {
		return b.Palette.OffColor
	}
	return b.Palette.OnColor
}

// BitAt returns the Bit color of the pixel at (x, y).
// BitAt(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// BitAt(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (b *Binary) BitAt(x, y int) Bit {
	if !(image.Point{x, y}.In(b.Rect)) {
		return Off
	}
	i := b.PixOffset(x, y)
	if b.Pix[i] == 0x0 {
		return Off
	}
	return On
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (b *Binary) PixOffset(x, y int) int {
	return (y-b.Rect.Min.Y)*b.Stride + (x-b.Rect.Min.X)*1
}

// Set sets the color of the pixel at (x, y).
//
// c is converted to Bit using the Binary image color model.
func (b *Binary) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = uint8(b.Palette.ConvertBit(c))
}

// SetBit sets the Bit of the pixel at (x, y).
func (b *Binary) SetBit(x, y int, c Bit) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = uint8(c)
}

// SetRect sets all the pixels in the rectangle defined by given rectangle.
func (b *Binary) SetRect(r image.Rectangle, c Bit) {
	r = r.Intersect(b.Rect)
	if !r.Empty() {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			i := b.PixOffset(r.Min.X, y)
			j := b.PixOffset(r.Max.X, y)
			// loop on all pixels (bytes) of this horizontal line
			for x := i; x < j; x++ {
				b.Pix[x] = uint8(c)
			}
		}
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (b *Binary) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(b.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Binary{}
	}
	i := b.PixOffset(r.Min.X, r.Min.Y)
	return &Binary{
		Pix:     b.Pix[i:],
		Stride:  b.Stride,
		Rect:    r,
		Palette: b.Palette,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (b *Binary) Opaque() bool {
	return true
}

// New returns a new Binary image with given width, height and binary palette.
func New(r image.Rectangle, p Palette) *Binary {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Binary{pix, 1 * w, r, p}
}

// NewFromImage converts src image into a Binary.Image with the given binary
// palette.
func NewFromImage(src image.Image, p Palette) *Binary {
	dst := New(src.Bounds(), p)
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}
