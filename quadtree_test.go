package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, resolution int) {
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
		_, err := NewBUQuadtree(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkBUQuadtreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 2)
}

func BenchmarkBUQuadtreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 4)
}

func BenchmarkBUQuadtreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 8)
}

func BenchmarkBUQuadtreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", 16)
}
