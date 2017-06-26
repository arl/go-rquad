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
	//scanner    imgscan.Scanner // reference image
}

// NewBasicTree creates a basic region quadtree from a scannable rectangular
// area and populates it with basic node instances.
//
// resolution is the smallest size in pixels that can have a leaf node, no
// further subdivisions will be performed on a node if its width or height is
// equal to this value.
//func OldNewBasicTree(scanner imgscan.Scanner, resolution int) (*BasicTree, error) {
//if resolution < 1 {
//return nil, errors.New("resolution must be greater than 0")
//}

//// To ensure a consistent behavior and eliminate corner cases,
//// the Quadtree's root node needs to have children. Thus, the
//// first instantiated Node needs to always be subdivided.
//// This condition asserts the resolution is respected.
//minDim := scanner.Bounds().Dx()
//if scanner.Bounds().Dy() < minDim {
//minDim = scanner.Bounds().Dy()
//}
//if minDim < resolution*2 {
//return nil, errors.New("the image smaller dimension must be greater or equal to twice the resolution")
//}

////root := &BasicNode{
////color:  Gray,
////bounds: scanner.Bounds(),
////}

//// create quadtree
//q := &BasicTree{
//resolution: resolution,
//scanner:    scanner,
//}
//// create root node
//q.root = q.nodeSetter.RootNode(scanner.bounds)

//q.subdivide(root)
//return q, nil
//}

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

	// BEGIN pseudo-code
	/*

		here with some well-chosen interfaces, we can scan the node, get the value
		and set the node value.

		The tree should be provided upon construction with a NodeValuer interface
		here we should call NodeValuer.SetValue

		Valuer.SetValue(node)

		func (v Valuer)SetValue(n Node) {
			// SetValue knows how to set the node value to something meaningful
			// get its color by scanning the image or whatever

			// so here for example if the node represents a region of  binary
			// colors, the node color will be set

			// setting the node color will, directly or indirectly, inform the node
			// wether it represents a leaf node or not, so that, later, calling
			// node.IsLeaf() will report the correct value.
		}

		if !n.IsLeaf() {
			q.subdivide(n)
		} else {
			// stop subdivision
		}

		// END pseudo-code

		uniform, col := q.scanner.IsUniform(bounds)
		switch uniform {
		case true:
			// quadrant is uniform, won't need to subdivide any further
			if col == binimg.White {
				n.color = White
			} else {
				n.color = Black
			}
		case false:
			// if we reached maximal resolution..
			if n.bounds.Dx()/2 < q.resolution || n.bounds.Dy()/2 < q.resolution {
				// ...make this node a black leaf, instead of gray
				n.color = Black
			} else {
				q.subdivide(n)
			}
		}
	*/

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
