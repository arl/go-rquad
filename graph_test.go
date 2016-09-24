package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func createScannerFromPNG(pngfile string) (binimg.Scanner, error) {
	var (
		bm      image.Image
		err     error
		scanner binimg.Scanner
	)
	if bm, err = loadPNG(pngfile); err != nil {
		return nil, err
	}
	if scanner, err = binimg.NewScanner(bm); err != nil {
		return nil, err
	}
	return scanner, nil
}

func createGraphFromScanner(scanner binimg.Scanner, resolution int) (*Graph, error) {
	q, err := NewBUQuadtree(scanner, resolution)
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
		createGraphFromScanner(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkGraphCreation(b *testing.B) {
	benchmarkGraphCreation(b, "./testdata/big.png", 4)
}
