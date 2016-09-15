package bmp

import "image"

// BruteForceScanner is the naive implementation of a bitmap scanner, it checks
// every pixel consecutively.
type BruteForceScanner struct {
	b *Bitmap
}

func (s *BruteForceScanner) SetBmp(bm *Bitmap) {
	s.b = bm
}

func (s BruteForceScanner) IsWhite(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x < bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != byte(White) {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s BruteForceScanner) IsBlack(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x < bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != byte(Black) {
				// immediately returns at the first 1 found
				return false
			}
		}
	}
	return true
}

func (s BruteForceScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	// naive implementation: check every pixel consecutively
	var yidx int

	// get first pixel color
	col := s.b.Bits[s.b.Width*topLeft.Y+topLeft.X]
	for y := topLeft.Y; y < bottomRight.Y; y++ {
		yidx = s.b.Width * y
		for x := topLeft.X; x < bottomRight.X; x++ {
			if s.b.Bits[x+yidx] != col {
				// immediately returns if color is different
				return Gray
			}
		}
	}
	return Color(col)
}
