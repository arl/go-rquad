package quadtree

import (
	"fmt"
	"testing"

	"github.com/aurelien-rainone/go-quadtrees/bmp"
)

func benchmarkQuadtreeCreation(b *testing.B, pngfile string, scanner bmp.Scanner, resolution int) {
	var (
		bm  *bmp.Bitmap
		err error
	)

	fmt.Println("loading", pngfile)
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

func BenchmarkBUQuadtreeCreationBruteForceScanner(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.BruteForceScanner{}, 4)
}

func BenchmarkBUQuadtreeCreationLinesScanner(b *testing.B) {
	benchmarkQuadtreeCreation(b, "./testdata/big.png", &bmp.LinesScanner{}, 4)
}
