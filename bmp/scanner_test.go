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
	testIsWhite(t, &bruteForceScanner{})
}

func TestBruteForceScannerIsBlack(t *testing.T) {
	testIsBlack(t, &bruteForceScanner{})
}

func TestBruteForceScannerIsFilled(t *testing.T) {
	testIsFilled(t, &bruteForceScanner{})
}

func BenchmarkBruteForceScanner(b *testing.B) {
	var (
		bm  *Bitmap
		err error
	)

	bm, err = loadPNG("./testdata/big.png")
	checkB(b, err)

	bm.SetScanner(&bruteForceScanner{})

	// run N times
	for n := 0; n < b.N; n++ {
		bm.IsWhite(image.Point{0, 0}, image.Point{bm.Width, bm.Height})
	}
}
