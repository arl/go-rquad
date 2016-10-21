package rquad

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/binimg"
)

// CNTree implements a Cardinal Neighbour Quadtree, a quadtree structure that
// allows finding neighbor quadrants in constant time O(1) regardless of their
// sizes.
//
// The time complexity reduction is obtained through the addition of four
// pointers per node in the quadtree, those pointers are the cardinal neighbour
// of the node.
//
// This quadtree structure has been proposed by Safwan W. Qasem, King Saud
// University, Kingdom of Saudi Arabia, in his paper "Cardinal Neighbor
// Quadtree: a New Quadtree-based Structure for Constant-Time Neighbor Finding"
type CNTree struct {
	resolution int            // maximal resolution
	scanner    binimg.Scanner // reference image
	root       *cnNode        // root node
	leaves     NodeList       // leaf nodes (filled during creation)
}

// NewCNTree creates a CNTree and populates it with cnNode instances,
// according to the content of the scanned image. It works only on square and
// power of 2 sized images, NewCNTree will return a nil CNTree pointer and an
// error if that's not the case.
//
// resolution is the minimal dimension of a leaf node, no further subdivisions
// will be performed on a leaf if its dimension is equal to the resolution.
func NewCNTree(scanner binimg.Scanner, resolution int) (*CNTree, error) {
	if !binimg.IsPowerOf2Image(scanner) {
		return nil, errors.New("image must be a square with power-of-2 dimensions")
	}

	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}

	// To ensure a consistent behavior and eliminate corner cases,
	// the Quadtree's root node needs to have children. Thus, the
	// first instantiated cnNode needs to always be subdivided.
	// This condition asserts the resolution is respected.
	if scanner.Bounds().Dx() < resolution*2 {
		return nil, errors.New("the image size must be greater or equal to twice the resolution")
	}

	q := &CNTree{
		resolution: resolution,
		scanner:    scanner,
	}

	q.root = q.newNode(q.scanner.Bounds(), nil, rootQuadrant)
	q.subdivide(q.root)
	return q, nil
}

func (q *CNTree) newNode(bounds image.Rectangle, parent *cnNode, location Quadrant) *cnNode {
	n := &cnNode{
		color:    Gray,
		bounds:   bounds,
		parent:   parent,
		location: location,
		size:     bounds.Dx(),
	}

	uniform, col := q.scanner.Uniform(bounds)
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
		if n.size/2 < q.resolution {
			// ...make this node a black leaf, instead of gray
			n.color = Black
		}
	}

	// fills leaves slices
	if n.color != Gray {
		q.leaves = append(q.leaves, n)
	}
	return n
}

func (q *CNTree) subdivide(p *cnNode) {
	// Step 1: Decomposing the gray quadrant and updating the
	//         parent node following the Z-order traversal.

	//     x0   x1     x2
	//  y0 .----.-------.
	//     |    |       |
	//     | NW |  NE   |
	//     |    |       |
	//  y1 '----'-------'
	//     | SW |  SE   |
	//  y2 '----'-------'
	//

	x0 := p.bounds.Min.X
	x1 := p.bounds.Min.X + p.size/2
	x2 := p.bounds.Max.X

	y0 := p.bounds.Min.Y
	y1 := p.bounds.Min.Y + p.size/2
	y2 := p.bounds.Max.Y

	// decompose current node in 4 sub-quadrants
	nw := q.newNode(image.Rect(x0, y0, x1, y1), p, Northwest)
	ne := q.newNode(image.Rect(x1, y0, x2, y1), p, Northeast)
	sw := q.newNode(image.Rect(x0, y1, x1, y2), p, Southwest)
	se := q.newNode(image.Rect(x1, y1, x2, y2), p, Southeast)

	// at creation, each sub-quadrant first inherits its parent external neighbours
	nw.cn[West] = p.cn[West]   // inherited
	nw.cn[North] = p.cn[North] // inherited
	nw.cn[East] = ne           // set for decomposition, will be updated after
	nw.cn[South] = sw          // set for decomposition, will be updated after
	ne.cn[West] = nw           // set for decomposition, will be updated after
	ne.cn[North] = p.cn[North] // inherited
	ne.cn[East] = p.cn[East]   // inherited
	ne.cn[South] = se          // set for decomposition, will be updated after
	sw.cn[West] = p.cn[West]   // inherited
	sw.cn[North] = nw          // set for decomposition, will be updated after
	sw.cn[East] = se           // set for decomposition, will be updated after
	sw.cn[South] = p.cn[South] // inherited
	se.cn[West] = sw           // set for decomposition, will be updated after
	se.cn[North] = ne          // set for decomposition, will be updated after
	se.cn[East] = p.cn[East]   // inherited
	se.cn[South] = p.cn[South] // inherited

	p.c[Northwest] = nw
	p.c[Northeast] = ne
	p.c[Southwest] = sw
	p.c[Southeast] = se

	p.updateNorthEast()
	p.updateSouthWest()

	// update all neighbours accordingly. After the decomposition
	// of a quadrant, all its neighbors in the four directions
	// must be informed of the change so that they can update
	// their own cardinal neighbors accordingly.
	p.updateNeighbours()

	// subdivide non-leaf nodes
	if nw.color == Gray {
		q.subdivide(nw)
	}
	if ne.color == Gray {
		q.subdivide(ne)
	}
	if sw.color == Gray {
		q.subdivide(sw)
	}
	if se.color == Gray {
		q.subdivide(se)
	}
}

// Root returns the quadtree root node.
func (q *CNTree) Root() Node {
	return q.root
}

// ForEachLeaf calls the given function for each leaf node of the quadtree.
//
// Successive calls to the provided function are performed in no particular
// order. The color parameter allows to loop on the leaves of a particular
// color, Black or White.
// NOTE: As by definition, Gray leaves do not exist, passing Gray to
// ForEachLeaf should return all leaves, independently of their color.
func (q *CNTree) ForEachLeaf(color Color, fn func(Node)) {
	for _, n := range q.leaves {
		if color == Gray || n.Color() == color {
			fn(n)
		}
	}
}
