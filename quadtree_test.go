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

func BenchmarkBUQuadtreeCreationCornerScannerRes2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.CornerScanner{}, 2)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.CornerScanner{}, 4)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.CornerScanner{}, 8)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.CornerScanner{}, 16)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes2MinWidth4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(4), 2)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes4MinWidth4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(4), 4)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes8MinWidth4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(4), 8)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes16MinWidth4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(4), 16)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes2MinWidth8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(8), 2)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes4MinWidth8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(8), 4)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes8MinWidth8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(8), 8)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes16MinWidth8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(8), 16)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes2MinWidth16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(16), 2)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes4MinWidth16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(16), 4)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes8MinWidth16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(16), 8)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes16MinWidth16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(16), 16)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes2MinWidth32(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(32), 2)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes4MinWidth32(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(32), 4)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes8MinWidth32(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(32), 8)
}

func BenchmarkBUQuadtreeCreationCornerScannerRes16MinWidth32(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewCornerScanner(32), 16)
}
