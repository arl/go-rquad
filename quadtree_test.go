package quadtree

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, fn newQuadtreeFunc, resolution int) {
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
		_, err := fn(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkBasicTreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 2)
}

func BenchmarkBasicTreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 4)
}

func BenchmarkBasicTreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 8)
}

func BenchmarkBasicTreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newBasicTree, 16)
}

func BenchmarkCNQuadtreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 2)
}

func BenchmarkCNQuadtreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 4)
}

func BenchmarkCNQuadtreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 8)
}

func BenchmarkCNQuadtreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/bigsquare.png", newCNQuadtree, 16)
}
