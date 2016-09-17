package bmp

import (
	"image"
	"testing"
)

func testIsWhite(t *testing.T, scanner Scanner) {
	ss := []string{
		"000",
		"100",
		"011",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
	}

	scanner.SetBmp(NewBitmapFromStrings(ss))

	for _, tt := range testTbl {
		actual := scanner.IsWhite(image.Point{tt.minx, tt.miny}, image.Point{tt.maxx, tt.maxy})
		if actual != tt.expected {
			t.Errorf("testIsWhite (%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsBlack(t *testing.T, scanner Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
		{2, 2, 3, 3, true},
	}

	scanner.SetBmp(NewBitmapFromStrings(ss))

	for _, tt := range testTbl {
		actual := scanner.IsBlack(image.Point{tt.minx, tt.miny}, image.Point{tt.maxx, tt.maxy})
		if actual != tt.expected {
			t.Errorf("testIsBlack (%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsFilled(t *testing.T, scanner Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}
	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               Color
	}{
		{0, 0, 3, 3, Gray},
		{1, 1, 3, 3, Gray},
		{0, 1, 1, 2, Black},
		{0, 0, 1, 1, White},
		{1, 0, 2, 1, White},
		{1, 0, 3, 2, White},
	}

	scanner.SetBmp(NewBitmapFromStrings(ss))

	for _, tt := range testTbl {
		actual := scanner.IsFilled(image.Point{tt.minx, tt.miny}, image.Point{tt.maxx, tt.maxy})
		if actual != tt.expected {
			t.Errorf("testIsFilled (%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func TestBruteForceScannerIsWhite(t *testing.T) {
	testIsWhite(t, &BruteForceScanner{})
}

func TestBruteForceScannerIsBlack(t *testing.T) {
	testIsBlack(t, &BruteForceScanner{})
}

func TestBruteForceScannerIsFilled(t *testing.T) {
	testIsFilled(t, &BruteForceScanner{})
}

func TestLinesScannerIsWhite(t *testing.T) {
	testIsWhite(t, &LinesScanner{})
}

func TestLinesScannerIsBlack(t *testing.T) {
	testIsBlack(t, &LinesScanner{})
}

func TestLinesScannerIsFilled(t *testing.T) {
	testIsFilled(t, &LinesScanner{})
}

func benchmarkScanner(b *testing.B, pngfile string, scanner Scanner) {
	var (
		bm  *Bitmap
		err error
	)

	bm, err = loadPNG(pngfile)
	checkB(b, err)

	bm.SetScanner(scanner)

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bm.IsWhite(image.Point{0, 0}, image.Point{bm.Width, bm.Height})
		bm.IsBlack(image.Point{0, 0}, image.Point{bm.Width, bm.Height})
		bm.IsFilled(image.Point{0, 0}, image.Point{bm.Width, bm.Height})
	}
}

func BenchmarkBruteForceScanner(b *testing.B) {
	benchmarkScanner(b, "./testdata/big.png", &BruteForceScanner{})
}

func BenchmarkLinesScanner(b *testing.B) {
	benchmarkScanner(b, "./testdata/big.png", &LinesScanner{})
}
