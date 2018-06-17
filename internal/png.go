package internal

import (
	"image"
	"image/png"
	"os"

	"github.com/arl/imgtools/binimg"
)

// LoadPNG is an helper function that uses binimg.NewFromImage internally.
func LoadPNG(filename string) (*binimg.Image, error) {
	var (
		f   *os.File
		img image.Image
		bm  *binimg.Image
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

// SavePNG is an helper function that encodes img to PNG and saves it on disk.
func SavePNG(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	return err
}
