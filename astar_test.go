package quadtree

import (
	"fmt"
	"image"
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
		float64(node.TopLeft().X), float64(node.TopLeft().Y),
		node.width(), node.height())
	ctx.Stroke()
}

func TestAStar(t *testing.T) {
	var (
		bm      image.Image
		g       *Graph
		err     error
		scanner binimg.Scanner
	)

	pngfile := "./testdata/big.png"
	resolution := 16

	bm, err = loadPNG(pngfile)
	check(t, err)
	scanner, err = binimg.NewScanner(bm)
	check(t, err)

	g, err = createGraphFromImage(scanner, resolution)
	check(t, err)

	fmt.Println("graph: nodes", len(g.nodes))
	org, dst := g.nodes[300], g.nodes[1300]

	ctx := gg.NewContextForImage(bm)
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
