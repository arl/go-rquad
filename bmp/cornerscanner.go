package bmp

import "image"

// CornerScanner checks the corners and other remarkable points of a
// rectangular region  before eventually scanning it completely.
//
// In the hope to avoid a full region scanning, some remarkable points will be
// checked before. Those remarkable points are:
// - the 4 corners: top-left, top-right, bottom-left and top-right.
// - the rectangle centre.
// - the 4 edges centres: top, right, bottom and left edges.
type CornerScanner struct {
	b     *Bitmap
	lines Scanner
}

func NewCornerScanner(maxBruteForceWidth int) Scanner {
	return &CornerScanner{}
}

func (s *CornerScanner) SetBmp(bm *Bitmap) {
	s.b = bm
	s.lines = &LinesScanner{}
	s.lines.SetBmp(bm)
}

// Returns true if all the remarkable pixels have the value of c.
func (s *CornerScanner) checkRemarkablePixels(topLeft, bottomRight image.Point, c byte) bool {
	bits := s.b.Bits
	x0, y0 := topLeft.X, topLeft.Y
	x2, y2 := bottomRight.X-1, topLeft.Y-1
	x1, y1 := (x0+x2)/2, (y0+y2)/2
	w := s.b.Width

	//     x0 x1 x2
	// y0  +--+--+
	// y1  +  +  +
	// y2  +--+--+

	if x0 == x2 || y0 == y0 {
		// this is flat rectangle, i.e one width or height is one, let the full
		// scan handle that
		return true
	}

	// top-left corner
	if bits[x0+w*y0] != c {
		return false
	}

	// top edge centre
	if bits[x1+w*y0] != c {
		return false
	}

	// top-right corner
	if bits[x2+w*y0] != c {
		return false
	}

	// left edge centre
	if bits[x0+w*y1] != c {
		return false
	}

	// rectange centre
	if bits[x1+w*y1] != c {
		return false
	}

	// left edge centre
	if bits[x2+w*y1] != c {
		return false
	}

	// bottom-left corner
	if bits[x0+w*y2] != c {
		return false
	}

	// bottom edge centre
	if bits[x1+w*y2] != c {
		return false
	}

	// bottom-right corner
	if bits[x2+w*y2] != c {
		return false
	}
	return true
}

func (s *CornerScanner) IsWhite(topLeft, bottomRight image.Point) bool {
	if !s.checkRemarkablePixels(topLeft, bottomRight, byte(White)) {
		return false
	}
	return s.lines.IsWhite(topLeft, bottomRight)
}

func (s *CornerScanner) IsBlack(topLeft, bottomRight image.Point) bool {
	if !s.checkRemarkablePixels(topLeft, bottomRight, byte(Black)) {
		return false
	}
	return s.lines.IsBlack(topLeft, bottomRight)
}

func (s *CornerScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	// color of top-left corner
	c := s.b.Bits[topLeft.X+topLeft.Y*s.b.Width]
	if !s.checkRemarkablePixels(topLeft, bottomRight, c) {
		return Gray
	}
	return s.lines.IsFilled(topLeft, bottomRight)
}
