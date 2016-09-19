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

func listNodes(n Quadnode) []Quadnode {
	var _listNodes func(n Quadnode, nodes *[]Quadnode)
	_listNodes = func(n Quadnode, nodes *[]Quadnode) {
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
	nodes := []Quadnode{}
	_listNodes(n, &nodes)
	return nodes
}
