package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func createGraphFromImage(scanner binimg.Scanner, res int) (*Graph, error) {
	q, err := NewBUQuadtree(scanner, res)
	if err != nil {
		return nil, err
	}
	return NewGraphFromQuadtree(q, nil), nil
}

func benchmarkGraphCreation(b *testing.B, pngfile string, resolution int) {
	var (
		bm      image.Image
		err     error
		scanner binimg.Scanner
	)

	bm, err = loadPNG(pngfile)
	checkB(b, err)
	scanner, err = binimg.NewScanner(bm)
	checkB(b, err)

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		createGraphFromImage(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkGraphCreation(b *testing.B) {
	benchmarkGraphCreation(b, "./testdata/big.png", 4)
}
