package rquad

import (
	"image"
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/binimg"
)

func benchmarkPointLocation(b *testing.B, fn newQuadtreeFunc, numPoints int) {
	var (
		img     *binimg.Binary
		scanner binimg.Scanner
		err     error
	)
	img, err = loadPNG("./testdata/bigsquare.png")
	checkB(b, err)

	r := rand.New(rand.NewSource(99))

	scanner, err = binimg.NewScanner(img)
	checkB(b, err)

	// create a quadtree
	q, err := fn(scanner, 8)
	checkB(b, err)

	randomPt := func(rect image.Rectangle) image.Point {
		return image.Pt(r.Intn(rect.Max.X-rect.Min.X)+rect.Min.X,
			r.Intn(rect.Max.Y-rect.Min.Y)+rect.Min.Y)
	}

	// fill a slice with random points
	points := make([]image.Point, numPoints, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = randomPt(q.Root().Bounds())
	}

	// run N times
	b.ResetTimer()
	var pt_ Node
	for n := 0; n < b.N; n++ {
		for _, pt := range points {
			// we don't want the compiler to optimize out the call to PointLocation
			// so we assign to a variable
			pt_ = PointLocation(q, pt)
		}
	}
	b.StopTimer()
	if pt_ == nil {
		b.Log("to not optimize the call to PointLocation, but should not be seen")
	}
}

func BenchmarkBasicQuadtreePointLocation10(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 10)
}

func BenchmarkBasicQuadtreePointLocation50(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 50)
}

func BenchmarkBasicQuadtreePointLocation200(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 200)
}

func BenchmarkBasicQuadtreePointLocation1000(b *testing.B) {
	benchmarkPointLocation(b, newBasicTree, 1000)
}

func BenchmarkCNTreePointLocation10(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 10)
}

func BenchmarkCNTreePointLocation50(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 50)
}

func BenchmarkCNTreePointLocation200(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 200)
}

func BenchmarkCNTreePointLocation1000(b *testing.B) {
	benchmarkPointLocation(b, newCNTree, 1000)
}
