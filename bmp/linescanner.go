package bmp

import (
	"bytes"
	"image"
)

// LinesScanner scans rectangular regions one line at at time.
type LinesScanner struct {
	b *Bitmap
}

func (s *LinesScanner) SetBmp(bm *Bitmap) {
	s.b = bm
}

func (s *LinesScanner) IsWhite(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		yidx = s.b.Width * y

		// look for the first Black byte of the line
		if bytes.IndexByte(s.b.Bits[yidx+topLeft.X:yidx+bottomRight.X], byte(Black)) != -1 {
			return false
		}
	}
	return true
}

func (s *LinesScanner) IsBlack(topLeft, bottomRight image.Point) bool {
	var yidx int

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		yidx = s.b.Width * y

		// look for the first White byte of the line
		if bytes.IndexByte(s.b.Bits[yidx+topLeft.X:yidx+bottomRight.X], byte(White)) != -1 {
			return false
		}
	}
	return true
}

func (s *LinesScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	var (
		corner Color // color of first pixel found
		target byte  // color we'll be looking for
	)

	yidx := s.b.Width * topLeft.Y
	corner = Color(s.b.Bits[yidx+topLeft.X])
	if corner == White {
		target = byte(Black)
	} else {
		target = byte(White)
	}

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		// look for the first byte of another color than the corner pixel
		if bytes.IndexByte(s.b.Bits[yidx+topLeft.X:yidx+bottomRight.X], target) != -1 {
			return Gray
		}
		yidx += s.b.Width
	}
	return corner
}
