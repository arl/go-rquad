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

func listNodes(n QNode) QNodeList {
	var _listNodes func(n QNode, nodes *QNodeList)
	_listNodes = func(n QNode, nodes *QNodeList) {
		switch n.Color() {
		case Gray:
			_listNodes(n.NorthWest(), nodes)
			_listNodes(n.NorthEast(), nodes)
			_listNodes(n.SouthWest(), nodes)
			_listNodes(n.SouthEast(), nodes)
		case White:
			*nodes = append(*nodes, n)
		}
	}
	nodes := QNodeList{}
	_listNodes(n, &nodes)
	return nodes
}
