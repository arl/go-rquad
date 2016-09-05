package bmp

import "image"

// Scanner is the interface implemented by objects that can scan a
// rectangular area of a Bitmap to check if its color is homogeneous
type Scanner interface {

	// IsWhite checks if all pixels of the given region are White.
	IsWhite(topLeft, bottomRight image.Point) bool

	// IsBlack checks if all pixels of the given region are Black.
	IsBlack(topLeft, bottomRight image.Point) bool

	// IsFilled checks if all pixels of the given region are of the same color.
	IsFilled(topLeft, bottomRight image.Point) Color

	// SetBmp defines the Bitmap on which the scanner performs the scan.
	SetBmp(bm *Bitmap)
}

// NaiveScanner is the naive implementation of a bitmap scanner, it checks
// every pixel consecutively.
type NaiveScanner struct {
	b *Bitmap
}

// NewNaiveScanner creates a new scanner, naively scanning the given bitmap
func NewNaiveScanner(b *Bitmap) *NaiveScanner {
	return &NaiveScanner{b}
}

func (s *NaiveScanner) SetBmp(bm *Bitmap) {
	s.b = bm
}

func (s NaiveScanner) IsWhite(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != White {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s NaiveScanner) IsBlack(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != Black {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s NaiveScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	// naive implementation: check every pixel consecutively
	var yidx int

	// get first pixel color
	col := s.b.Bits[s.b.Width*topLeft.Y+topLeft.X]
	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != col {
				// immediately returns if color is different
				return Gray
			}
		}
	}
	return col
}
