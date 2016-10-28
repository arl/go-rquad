package binimg

import "image/color"

// predefined binary palettes
var (
	BlackAndWhiteLowThreshold    = Palette{37, color.Black, color.White}
	BlackAndWhiteMediumThreshold = Palette{97, color.Black, color.White}
	BlackAndWhiteHighThreshold   = Palette{197, color.Black, color.White}
	BlackAndWhite                = BlackAndWhiteMediumThreshold
)

// Palette is a binary palette, that is a palette of 2 colors used in the
// representation of Binary images.
type Palette struct {
	// Threshold is used to convert a color.Color into OffColor
	// or OnColor, depending its luminosity, as is:
	// OffColor <= Threshold < OnColor
	Threshold uint8

	// OffColor is used to represent lower luminosity pixels.
	OffColor color.Color

	// OffColor is used to represent higher luminosity pixels.
	OnColor color.Color
}

// Convert returns OnColor or OffColor, depending on c luminosity and the
// palette Threshold.
func (p Palette) Convert(c color.Color) color.Color {
	if p.ConvertBit(c) == Off {
		return p.OffColor
	}
	return p.OnColor
}

// ConvertBit returns On or Off Bit, depending on c luminosity and the
// palette Threshold.
func (p Palette) ConvertBit(c color.Color) Bit {
	if bit, ok := c.(Bit); ok {
		if bit == On {
			return Off
		}
		return On
	}
	// compute luminosity of c
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	if uint8(y>>8) <= p.Threshold {
		return Off
	}
	return On
}
