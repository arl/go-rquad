package bmp

import "image"

// SmartScanner implements various scanning techniques and switch between them
// depending on the dimensions of the regions to scan.
type SmartScanner struct {
	b *Bitmap

	lines Scanner
	brute Scanner

	// max region width for which we still use brute force scanner
	maxBruteForceWidth int
}

func NewSmartScanner(maxBruteForceWidth int) Scanner {
	return &SmartScanner{
		maxBruteForceWidth: maxBruteForceWidth,
	}
}

func (s *SmartScanner) SetBmp(bm *Bitmap) {
	s.b = bm
	s.lines = &LinesScanner{}
	s.lines.SetBmp(bm)
	s.brute = &BruteForceScanner{}
	s.brute.SetBmp(bm)
}

func (s *SmartScanner) selectScanner(topLeft, bottomRight image.Point) Scanner {
	if bottomRight.X-topLeft.X < s.maxBruteForceWidth {
		// region is small enough, use brute force scanner
		return s.brute
	}
	return s.lines
}

func (s *SmartScanner) IsWhite(topLeft, bottomRight image.Point) bool {
	return s.selectScanner(topLeft, bottomRight).IsWhite(topLeft, bottomRight)
}

func (s *SmartScanner) IsBlack(topLeft, bottomRight image.Point) bool {
	return s.selectScanner(topLeft, bottomRight).IsBlack(topLeft, bottomRight)
}

func (s *SmartScanner) IsFilled(topLeft, bottomRight image.Point) Color {
	return s.selectScanner(topLeft, bottomRight).IsFilled(topLeft, bottomRight)
}
