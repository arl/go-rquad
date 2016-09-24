package quadtree

import (
	"testing"

	"github.com/aurelien-rainone/binimg"
	astar "github.com/beefsack/go-astar"
	"github.com/fogleman/gg"
)

func drawPath(ctx *gg.Context, path []astar.Pather) {
	ctx.SetRGB(0, 1, 0)
	ctx.SetLineWidth(4)
	for _, pt := range path {
		x, y := pt.(*Node).center()
		ctx.LineTo(x, y)
	}
	ctx.Stroke()
}

func drawNode(ctx *gg.Context, node *Node) {
	ctx.SetRGB(1, 0, 0)
	ctx.SetLineWidth(1)
	ctx.DrawRectangle(
		float64(node.Bounds().Min.X), float64(node.Bounds().Min.Y),
		node.width(), node.height())
	ctx.Stroke()
}

func TestAStar(t *testing.T) {
	var (
		g       *Graph
		err     error
		scanner binimg.Scanner
	)

	pngfile := "./testdata/big.png"
	resolution := 16

	scanner, err = createScannerFromPNG(pngfile)
	check(t, err)

	g, err = createGraphFromScanner(scanner, resolution)
	check(t, err)

	org, dst := g.nodes[300], g.nodes[1300]

	ctx := gg.NewContextForImage(scanner)
	drawNode(ctx, org)
	drawNode(ctx, dst)

	path, _, found := astar.Path(org, dst)
	if found {
		for _, p := range path {
			drawNode(ctx, p.(*Node))
		}
		drawPath(ctx, path)
	} else {
		t.Errorf("path not found")
	}
	savePNG(ctx.Image(), "testpath.png")
}

func benchmarkAStar(b *testing.B, pngfile string, resolution int, orgidx, dstidx int) {

	var (
		g       *Graph
		err     error
		scanner binimg.Scanner
	)

	scanner, err = createScannerFromPNG(pngfile)
	checkB(b, err)

	g, err = createGraphFromScanner(scanner, resolution)
	checkB(b, err)

	var found bool

	b.ResetTimer()
	// run N times
	for n := 0; n < b.N; n++ {
		if _, _, found = astar.Path(g.nodes[orgidx], g.nodes[dstidx]); !found {
			b.Errorf("path should have been found")
		}
	}
}

func BenchmarkAStar(b *testing.B) {
	benchmarkAStar(b, "./testdata/big.png", 16, 300, 1300)
}
