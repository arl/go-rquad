// Package binimg proposes an in-memory binary image format, implementing the
// image.Image interface, alongside a set of efficient tools to scan
// rectangular regions of such images. A binary image has only two possible
// colors for each pixel, generally Black and White, though any two colors can
// be used.
//
// Though the information represented by each pixel could be stored as a single
// bit, and thus take a smaller memory footprint, choice has been made to
// represent Bit pixels as byte values, that can either be 0 (Black or Off) or
// 255 (White or On), mostly for simplicity reasons.
//
// Binary images are created either by calling functions such as NewFromImage
// and NewBinary, or their counterparts accepting a custom binaryModel.
package binimg

import (
	"image"
	"image/color"
	"image/draw"
)

// Black and White are the only colors that a Binary image pixel can have.
var (
	Black = Bit{0}
	White = Bit{255}
)

// Alias colors for Black and White
var (
	Off = Black
	On  = White
)

// Bit represents a Black or White only binary color.
type Bit struct {
	v byte
}

// RGBA returns the red, green, blue and alpha values for a Bit color.
//
// alpha is always 0xffff (fully opaque) and r, g, b are all 0 or all 0xffff.
func (c Bit) RGBA() (r, g, b, a uint32) {
	v := uint32(c.v)
	v |= v << 8
	return v, v, v, 0xffff
}

// Other returns a Bit with the other value.
func (c Bit) Other() Bit {
	if c.v == 0 {
		return White
	}
	return Black
}

// Various binary models with different thresholds.
var (
	BinaryModelLowThreshold    = NewBinaryModel(37)
	BinaryModelMediumThreshold = NewBinaryModel(97)
	BinaryModelHighThreshold   = NewBinaryModel(197)
	BinaryModel                = BinaryModelMediumThreshold
)

type binaryModel struct {
	threshold uint8
}

func (m binaryModel) Convert(c color.Color) color.Color {
	if _, ok := c.(Bit); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	y := (299*r + 587*g + 114*b + 500) / 1000
	if uint8(y>>8) > m.threshold {
		return White
	}
	return Black
}

// NewBinaryModel creates a new binaryModel that converts any color to a Bit.
//
// binaryModel is an opaque (as in not exported) type. The threshold is the
// limit over which source colors are converted to White, under to Black.
func NewBinaryModel(threshold uint8) binaryModel {
	return binaryModel{threshold}
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

	model binaryModel
}

// ColorModel returns the image.Image's color model.
func (b *Binary) ColorModel() color.Model { return b.model }

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (b *Binary) Bounds() image.Rectangle { return b.Rect }

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (b *Binary) At(x, y int) color.Color {
	return b.BitAt(x, y)
}

// BitAt returns the Bit color of the pixel at (x, y).
// BitAt(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the
// grid. BitAt(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right
// one.
func (b *Binary) BitAt(x, y int) Bit {
	if !(image.Point{x, y}.In(b.Rect)) {
		return Bit{}
	}
	i := b.PixOffset(x, y)
	return Bit{b.Pix[i]}
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
	b.Pix[i] = b.model.Convert(c).(Bit).v
}

// SetBit sets the Bit of the pixel at (x, y).
func (b *Binary) SetBit(x, y int, c Bit) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = c.v
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
		Pix:    b.Pix[i:],
		Stride: b.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (b *Binary) Opaque() bool {
	return true
}

// New returns a new Binary image with the given bounds.
func New(r image.Rectangle) *Binary {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Binary{pix, 1 * w, r, BinaryModel}
}

// NewCustomBinary returns a new Binary image with the given bounds and binary
// model.
func NewCustomBinary(r image.Rectangle, model binaryModel) *Binary {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Binary{pix, 1 * w, r, model}
}

// NewFromImage returns the binary image that is the conversion of the given
// source image.
func NewFromImage(src image.Image) *Binary {
	dst := New(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}

// NewCustomFromImage returns the binary image that is the conversion of the
// given source image with the specified binary model.
func NewCustomFromImage(src image.Image, model binaryModel) *Binary {
	dst := NewCustomBinary(src.Bounds(), model)
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}
