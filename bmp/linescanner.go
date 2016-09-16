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
		if bytes.IndexByte(s.b.Bits[yidx:yidx+bottomRight.X], byte(Black)) != -1 {
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
		if bytes.IndexByte(s.b.Bits[yidx:yidx+bottomRight.X], byte(White)) != -1 {
			return false
		}
	}
	return true
}

func (s *LinesScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	var yidx int

	var (
		corner Color // color of first pixel found
		target byte  // color we'll be looking for
	)

	corner = Color(s.b.Bits[0])
	if corner == White {
		target = byte(Black)
	} else {
		target = byte(White)
	}

	for y := topLeft.Y; y < bottomRight.Y-1; y++ {
		yidx = s.b.Width * y

		// look for the first byte that is different than the corner pixel
		if bytes.IndexByte(s.b.Bits[yidx:yidx+bottomRight.X], target) != -1 {
			return Gray
		}
	}
	return corner
}
