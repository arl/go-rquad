package rquad

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/imgscan"
)

// BasicTree is a standard implementation of a region quadtree.
//
// It performs a standard quadtree subdivision of the rectangular area
// represented by an binimg.Scanner.
type BasicTree struct {
	resolution int             // leaf node resolution
	scanner    imgscan.Scanner // reference image
	root       Node            // root node
	leaves     NodeList        // leaf nodes (filled during creation)
}

// NewBasicTree creates a basic region quadtree from a scannable rectangular
// area and populates it with basic node instances.
//
// resolution is the smallest size in pixels that can have a leaf node, no
// further subdivisions will be performed on a node if its width or height is
// equal to this value.
func NewBasicTree(scanner imgscan.Scanner, resolution int) (*BasicTree, error) {
	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}

	// To ensure a consistent behavior and eliminate corner cases,
	// the Quadtree's root node needs to have children. Thus, the
	// first instantiated Node needs to always be subdivided.
	// This condition asserts the resolution is respected.
	minDim := scanner.Bounds().Dx()
	if scanner.Bounds().Dy() < minDim {
		minDim = scanner.Bounds().Dy()
	}
	if minDim < resolution*2 {
		return nil, errors.New("the image smaller dimension must be greater or equal to twice the resolution")
	}

	// create root node
	root := &basicNode{
		color:  Gray,
		bounds: scanner.Bounds(),
	}

	// create quadtree
	q := &BasicTree{
		resolution: resolution,
		scanner:    scanner,
		root:       root,
	}
	q.subdivide(root)
	return q, nil
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
		if color == Gray || n.Color() == color {
			fn(n)
		}
	}
}

func (q *BasicTree) newChildNode(bounds image.Rectangle, parent *basicNode, location Quadrant) *basicNode {
	n := &basicNode{
		color:    Gray,
		bounds:   bounds,
		parent:   parent,
		location: location,
	}

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

	// fills leaves slices
	if n.color != Gray {
		q.leaves = append(q.leaves, n)
	}
	return n
}

func (q *BasicTree) subdivide(n *basicNode) {
	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'
	//
	x0 := n.bounds.Min.X
	x1 := n.bounds.Min.X + n.bounds.Dx()/2
	x2 := n.bounds.Max.X

	y0 := n.bounds.Min.Y
	y1 := n.bounds.Min.Y + n.bounds.Dy()/2
	y2 := n.bounds.Max.Y

	// create the 4 children nodes, one per quadrant
	n.c[Northwest] = q.newChildNode(image.Rect(x0, y0, x1, y1), n, Northwest)
	n.c[Southwest] = q.newChildNode(image.Rect(x0, y1, x1, y2), n, Southwest)
	n.c[Northeast] = q.newChildNode(image.Rect(x1, y0, x2, y1), n, Northeast)
	n.c[Southeast] = q.newChildNode(image.Rect(x1, y1, x2, y2), n, Southeast)
}

// Root returns the quadtree root node.
func (q *BasicTree) Root() Node {
	return q.root
}

// basicNode represents a standard quadtree node.
//
// It is a basic implementation of the Node interface, the one used in the
// BasicTree implementation of the Quadtree interface.
type basicNode struct {
	parent   Node            // pointer to the parent node
	c        [4]Node         // children nodes
	bounds   image.Rectangle // node bounds
	color    Color           // node color
	location Quadrant        // node location inside its parent
}

// Parent returns the quadtree node that is the parent of current one.
func (n *basicNode) Parent() Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}

// Child returns current node child at specified quadrant.
func (n *basicNode) Child(q Quadrant) Node {
	if n.c[q] == nil {
		return nil
	}
	return n.c[q]
}

// Bounds returns the bounds of the rectangular area represented by this
// quadtree node.
func (n *basicNode) Bounds() image.Rectangle {
	return n.bounds
}

// Color returns the node Color.
func (n *basicNode) Color() Color {
	return n.color
}

// Location returns the node inside its parent quadrant
func (n *basicNode) Location() Quadrant {
	return n.location
}
