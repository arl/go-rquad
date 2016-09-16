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

func benchmarkBUQuadtreeCreationSmartScanner(b *testing.B, resolution, maxBruteForceWidth int) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(maxBruteForceWidth), resolution)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW8Res2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(8), 2)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW8Res4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(8), 4)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW8Res8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(8), 8)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW8Res16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(8), 16)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW16Res2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(16), 2)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW16Res4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(16), 4)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW16Res8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(16), 8)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW16Res16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(16), 16)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW32Res2(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(32), 2)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW32Res4(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(32), 4)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW32Res8(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(32), 8)
}

func BenchmarkBUQuadtreeCreationSmartScannerMaxBFW32Res16(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", bmp.NewSmartScanner(32), 16)
}
