// +build ignore
package bmp

import (
	"fmt"
)

func NewBitmapFromStrings(ss []string) *Bitmap {
	w, h := len(ss[0]), len(ss)
	for i := range ss {
		if len(ss[i]) != w {
			panic("all strings should have the same length")
		}
	}

	bmp := Bitmap{
		Width:  w,
		Height: h,
		Bits:   make([]Color, w*h),
	}

	for y := range ss {
		for x := range ss[y] {
			if ss[y][x] == '1' {
				bmp.Bits[x+w*y] = White
			}
		}
	}
	return &bmp
}

func (bmp Bitmap) String() string {
	var s string
	for y := 0; y < bmp.Height; y++ {
		for x := 0; x < bmp.Width; x++ {
			s += fmt.Sprintf("%d", bmp.Bits[x+bmp.Width*y])
		}
		s += "\n"
	}
	return s
}
