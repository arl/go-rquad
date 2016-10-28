package rquad

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
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

	bm = binimg.NewFromImage(img, binimg.BlackAndWhite)
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

type newQuadtreeFunc func(imgscan.Scanner, int) (Quadtree, error)

func newBasicTree(scanner imgscan.Scanner, resolution int) (Quadtree, error) {
	return NewBasicTree(scanner, resolution)
}

func newCNTree(scanner imgscan.Scanner, resolution int) (Quadtree, error) {
	return NewCNTree(scanner, resolution)
}

func appendNode(nl *NodeList) func(Node) {
	return func(n Node) {
		*nl = append(*nl, n)
	}
}

func neighbourColors(n Node) (white, black int) {
	ForEachNeighbour(n, func(nb Node) {
		switch nb.Color() {
		case Black:
			black++
		case White:
			white++
		}
	})
	return
}
