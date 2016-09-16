package quadtree

import (
	"testing"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, scanner bmp.Scanner, resolution int) {
	var (
		bm  *bmp.Bitmap
		err error
	)

	bm, err = loadPNG(pngfile)
	checkB(b, err)

	bm.SetScanner(scanner)

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := NewBUQuadtree(bm, resolution)
		checkB(b, err)
	}
}

func BenchmarkBUQuadtreeCreationBruteForceScannerRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.BruteForceScanner{}, 2)
}

func BenchmarkBUQuadtreeCreationBruteForceScannerRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.BruteForceScanner{}, 4)
}

func BenchmarkBUQuadtreeCreationBruteForceScannerRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.BruteForceScanner{}, 8)
}

func BenchmarkBUQuadtreeCreationBruteForceScannerRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.BruteForceScanner{}, 16)
}

func BenchmarkBUQuadtreeCreationLinesScannerRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.LinesScanner{}, 2)
}

func BenchmarkBUQuadtreeCreationLinesScannerRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.LinesScanner{}, 4)
}

func BenchmarkBUQuadtreeCreationLinesScannerRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.LinesScanner{}, 8)
}

func BenchmarkBUQuadtreeCreationLinesScannerRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.LinesScanner{}, 16)
}
