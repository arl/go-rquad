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
// binimg.Image images are created either by calling functions such as
// New and NewFromImage.
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
type Bit struct{ V byte }

// RGBA returns the red, green, blue and alpha values for a Bit color.
//
// alpha is always 0xffff (fully opaque) and r, g, b are all 0 or all 0xffff.
func (c Bit) RGBA() (r, g, b, a uint32) {
	v := uint32(c.V)
	v |= v << 8
	return v, v, v, 0xffff
}

// Other returns a Bit with the other value.
func (c Bit) Other() Bit {
	if c.V == 0 {
		return White
	}
	return Black
}

// Model is the color model for binary images.
var Model color.Model = color.ModelFunc(model)

func model(c color.Color) color.Color {
	if _, ok := c.(Bit); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	y := (299*r + 587*g + 114*b + 500) / 1000
	if uint8(y>>8) > 97 {
		return White
	}
	return Black
}

// Image is an in-memory image whose At method returns Bit values.
type Image struct {
	// Pix holds the image's pixels, as 0 or 1 uint8 values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the image.Image's color model.
func (b *Image) ColorModel() color.Model { return Model }

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (b *Image) Bounds() image.Rectangle { return b.Rect }

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (b *Image) At(x, y int) color.Color {
	return b.BitAt(x, y)
}

// BitAt returns the Bit color of the pixel at (x, y).
// BitAt(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the
// grid. BitAt(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right
// one.
func (b *Image) BitAt(x, y int) Bit {
	if !(image.Point{x, y}.In(b.Rect)) {
		return Bit{}
	}
	i := b.PixOffset(x, y)
	return Bit{b.Pix[i]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (b *Image) PixOffset(x, y int) int {
	return (y-b.Rect.Min.Y)*b.Stride + (x-b.Rect.Min.X)*1
}

// Set sets the color of the pixel at (x, y).
//
// c is converted to Bit using the image color model.
func (b *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = Model.Convert(c).(Bit).V
}

// SetBit sets the Bit of the pixel at (x, y).
func (b *Image) SetBit(x, y int, c Bit) {
	if !(image.Point{x, y}.In(b.Rect)) {
		return
	}
	i := b.PixOffset(x, y)
	b.Pix[i] = c.V
}

// SetRect sets all the pixels in the rectangle defined by given rectangle.
func (b *Image) SetRect(r image.Rectangle, c Bit) {
	r = r.Intersect(b.Rect)
	if !r.Empty() {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			i := b.PixOffset(r.Min.X, y)
			j := b.PixOffset(r.Max.X, y)
			// loop on all pixels (bytes) of this horizontal line
			for x := i; x < j; x++ {
				b.Pix[x] = c.V
			}
		}
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (b *Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(b.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Image{}
	}
	i := b.PixOffset(r.Min.X, r.Min.Y)
	return &Image{
		Pix:    b.Pix[i:],
		Stride: b.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (b *Image) Opaque() bool {
	return true
}

// New returns a new binary image with the given bounds.
func New(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	return &Image{pix, 1 * w, r}
}

// NewFromImage returns a new binary image that is the conversion of src image.
func NewFromImage(src image.Image) *Image {
	dst := New(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}
