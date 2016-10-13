package quadtree

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

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

// helper function that uses binimg.NewFromImage internally.
func loadPNG(filename string) (*binimg.Binary, error) {
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

func savePNG(img image.Image, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return nil
}

type newQuadtreeFunc func(binimg.Scanner, int) (Quadtree, error)

func newBUQuadtree(scanner binimg.Scanner, resolution int) (Quadtree, error) {
	return NewBUQuadtree(scanner, resolution)
}

func newCNQuadtree(scanner binimg.Scanner, resolution int) (Quadtree, error) {
	return NewCNQuadtree(scanner, resolution)
}

func appendNode(nl *QNodeList) func(QNode) {
	return func(n QNode) {
		*nl = append(*nl, n)
	}
}

func neighbourColors(n QNode) (white, black int) {
	var nodes QNodeList
	n.Neighbours(&nodes)
	for _, nb := range nodes {
		switch nb.Color() {
		case Black:
			black++
		case White:
			white++
		}
	}
	return
}
