package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

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
		q, err := NewBUQuadtree(scanner, resolution)
		checkB(b, err)
		NewGraphFromQuadtree(q)
	}
}

func BenchmarkGraphCreation(b *testing.B) {
	benchmarkGraphCreation(b, "./testdata/big.png", 4)
}
