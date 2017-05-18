package internal

import (
	"image"
	"image/png"
	"os"

	"github.com/aurelien-rainone/binimg"
)

// helper function that uses binimg.NewFromImage internally.
func LoadPNG(filename string) (*binimg.Binary, error) {
	var (
		f   *os.File
		img image.Image
		bm  *binimg.Binary
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

	bm = binimg.NewFromImage(img)
	return bm, nil
}

func SavePNG(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	return err
}
