package rquad

import (
	"image"
	"testing"

	"github.com/aurelien-rainone/go-rquad/internal"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

func benchmarkQuadtreeCreation(b *testing.B, fn newQuadtreeFunc, resolution int) {
	var (
		bm      image.Image
		err     error
		scanner imgscan.Scanner
	)

	bm, err = internal.LoadPNG("./testdata/bigsquare.png")
	checkB(b, err)
	scanner, err = imgscan.NewScanner(bm)
	checkB(b, err)

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := fn(scanner, resolution)
		checkB(b, err)
	}
}

func BenchmarkBasicCreationRes32(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 32)
}

func BenchmarkBasicCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 16)
}

func BenchmarkBasicCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 8)
}

func BenchmarkBasicCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 4)
}

func BenchmarkBasicCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 2)
}

func BenchmarkBasicCreationRes1(b *testing.B) {
	benchmarkQuadtreeCreation(b, newBasicTree, 1)
}

func BenchmarkCNTreeCreationRes32(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 32)
}

func BenchmarkCNTreeCreationRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 16)
}

func BenchmarkCNTreeCreationRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 8)
}

func BenchmarkCNTreeCreationRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 4)
}

func BenchmarkCNTreeCreationRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 2)
}

func BenchmarkCNTreeCreationRes1(b *testing.B) {
	benchmarkQuadtreeCreation(b, newCNTree, 1)
}
