package bmp

import "image"

// Scanner is the interface implemented by objects that can scan a
// rectangular area of a Bitmap to check if its color is homogeneous
type Scanner interface {

	// IsWhite checks if all elements
	IsWhite(b *Bitmap, topLeft, bottomRight image.Point) bool

	IsBlack(b *Bitmap, topLeft, bottomRight image.Point) bool

	IsFilled(b *Bitmap, topLeft, bottomRight image.Point) Color
}

// NaiveScanner is the naive implementation of a bitmap scanner, it checks
// every pixel consecutively.
type NaiveScanner struct{}

func (s NaiveScanner) IsWhite(b *Bitmap, topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if b.Bits[x+yidx] != White {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s NaiveScanner) IsBlack(b *Bitmap, topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if b.Bits[x+yidx] != Black {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s NaiveScanner) IsFilled(b *Bitmap, topLeft, bottomRight image.Point) Color {
	// naive implementation: check every pixel consecutively
	var yidx int

	// get first pixel color
	col := b.Bits[b.Width*topLeft.Y+topLeft.X]
	for y := topLeft.Y; y <= bottomRight.Y; y++ {
		yidx = b.Width * y
		for x := topLeft.X; x <= bottomRight.X; x++ {
			if b.Bits[x+yidx] != col {
				// immediately returns if color is different
				return Gray
			}
		}
	}
	return col
}
