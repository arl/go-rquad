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
