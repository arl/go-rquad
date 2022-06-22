package rquad

import (
	"image"
	"math/rand"
	"testing"

	"github.com/arl/go-rquad/internal"
	"github.com/arl/imgtools/binimg"
	"github.com/arl/imgtools/imgscan"
)

func benchmarkPointLocation(b *testing.B, fn newQuadtreeFunc, numPoints int, resolution int) {
	var (
		img     *binimg.Image
		scanner imgscan.Scanner
		err     error
	)
	img, err = internal.LoadPNG("./testdata/random-1024x1024.png")
	checkB(b, err)

	r := rand.New(rand.NewSource(99))

	scanner, err = imgscan.NewScanner(img)
	checkB(b, err)

	// create a quadtree
	q, err := fn(scanner, resolution)
	checkB(b, err)

	randomPt := func(rect image.Rectangle) image.Point {
		return image.Pt(r.Intn(rect.Max.X-rect.Min.X)+rect.Min.X,
			r.Intn(rect.Max.Y-rect.Min.Y)+rect.Min.Y)
	}

	// fill a slice with random points
	points := make([]image.Point, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = randomPt(q.Root().Bounds())
	}

	// run N times
	b.ResetTimer()
	var dummy Node
	for n := 0; n < b.N; n++ {
		for _, pt := range points {
			// we don't want the compiler to optimize out the call to PointLocation
			// so we assign to a variable
			dummy = Locate(q, pt)
		}
	}
	b.StopTimer()
	if dummy == nil {
		b.Log("to not optimize the call to PointLocation, but should not be seen")
	}
}

func BenchmarkBasicPointLocationRes32(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 32)
}

func BenchmarkBasicPointLocationRes16(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 16)
}

func BenchmarkBasicPointLocationRes8(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 8)
}

func BenchmarkBasicPointLocationRes4(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 4)
}

func BenchmarkBasicPointLocationRes2(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 2)
}

func BenchmarkBasicPointLocationRes1(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 100, 1)
}

func BenchmarkCNTreePointLocationRes32(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 32)
}

func BenchmarkCNTreePointLocationRes16(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 16)
}

func BenchmarkCNTreePointLocationRes8(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 8)
}

func BenchmarkCNTreePointLocationRes4(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 4)
}

func BenchmarkCNTreePointLocationRes2(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 2)
}

func BenchmarkCNTreePointLocationRes1(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 100, 1)
}
