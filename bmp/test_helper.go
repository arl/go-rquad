package bmp

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
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
		Bits:   make([]byte, w*h),
	}

	for y := range ss {
		for x := range ss[y] {
			if ss[y][x] == '1' {
				bmp.Bits[x+w*y] = byte(White)
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
func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Fatal(err)
	}
}

// helper function that uses NewFromImage internally.
func loadPNG(filename string) (*Bitmap, error) {
	var (
		f   *os.File
		img image.Image
		bm  *Bitmap
		err error
	)

	f, err = os.Open(filename)
	if err != nil {
		return bm, err
	}
	defer f.Close()

	img, err = png.Decode(f)
	if err != nil {
		return bm, err
	}

	bm = NewFromImage(img)
	return bm, nil
}
