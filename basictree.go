package rquad

import (
	"image"
)

// BasicTree is a standard implementation of a region quadtree.
//
// It performs a standard quadtree subdivision of the rectangular area
// represented by an binimg.Scanner.
type BasicTree struct {
	resolution int        // leaf node resolution
	nodeSetter NodeSetter // node setter
	root       Node       // root node
	leaves     NodeList   // leaf nodes (filled during creation)
}

func NewBasicTree(nodeSetter NodeSetter) *BasicTree {
	// create quadtree
	q := &BasicTree{
		nodeSetter: nodeSetter,
	}
	// create root node
	q.root = q.nodeSetter.NewRoot()

	q.subdivide(q.root)
	return q
}

// ForEachLeaf calls the given function for each leaf node of the quadtree.
//
// Successive calls to the provided function are performed in no particular
// order. The color parameter allows to loop on the leaves of a particular
// color, Black or White.
// NOTE: As by definition, Gray leaves do not exist, passing Gray to
// ForEachLeaf should return all leaves, independently of their color.
func (q *BasicTree) ForEachLeaf(color Color, fn func(Node)) {
	for _, n := range q.leaves {
		if n.IsLeaf() {
			fn(n)
		}
	}
}

func (q *BasicTree) newChildNode(bounds image.Rectangle, parent Node, location Quadrant) Node {
	n := q.nodeSetter.NewNode(parent, location, bounds)
	q.nodeSetter.ScanAndSet(&n)
	if n.IsLeaf() {
		// fills leaves slices
		q.leaves = append(q.leaves, n)
	} else {
		// nothing to do
		q.subdivide(n)
	}
	return n
}

func (q *BasicTree) subdivide(n Node) {
	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'
	//
	x0 := n.Bounds().Min.X
	x1 := n.Bounds().Min.X + n.Bounds().Dx()/2
	x2 := n.Bounds().Max.X

	y0 := n.Bounds().Min.Y
	y1 := n.Bounds().Min.Y + n.Bounds().Dy()/2
	y2 := n.Bounds().Max.Y

	// create the 4 children nodes, one per quadrant
	n.SetChild(Northwest, q.newChildNode(image.Rect(x0, y0, x1, y1), n, Northwest))
	n.SetChild(Southwest, q.newChildNode(image.Rect(x0, y1, x1, y2), n, Southwest))
	n.SetChild(Northeast, q.newChildNode(image.Rect(x1, y0, x2, y1), n, Northeast))
	n.SetChild(Southeast, q.newChildNode(image.Rect(x1, y1, x2, y2), n, Southeast))
}

// Root returns the quadtree root node.
func (q *BasicTree) Root() Node {
	return q.root
}
