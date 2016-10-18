package quadtree

import (
	"errors"
	"image"

	"github.com/aurelien-rainone/binimg"
)

// CNQuadtree is a quadtree based structure specically crafted for
// constant-time neighbour finding. It works on square and power-of-2 sized
// quadrants.
//
// The Cardinal Neighbor Quadtree, a pointer based data structure, can
// determine the existence, and access a smaller, equal or greater size
// neighbour in constant-time O(1). The time complexity reduction is obtained
// through the addition of only four pointers per leaf node in the quadtree.
//
// This quadtree structure has been proposed by Safwan W. Qasem, King Saud
// University, Kingdom of Saudi Arabia, in his paper "Cardinal Neighbor
// Quadtree: a New Quadtree-based Structure for Constant-Time Neighbor Finding"
type CNQuadtree struct {
	resolution int            // maximal resolution
	scanner    binimg.Scanner // reference image
	root       *CNQNode       // root node
	leaves     QNodeList      // leaf nodes (filled during creation)
}

// NewCNQuadtree creates a CNQuadtree and populates it with CNQNode's,
// according to the content of the scanned image. If the image is not a square
// having power-of-2 sides, the image will be redimensionned to fit this
// requirement.
//
// Resolution is the minimal dimension that can have a leaf node, no further
// subdivisions will be performed on a node if its width or height is equal to
// this value.
func NewCNQuadtree(scanner binimg.Scanner, resolution int) (*CNQuadtree, error) {
	// initialize package level variables
	initPackage()

	if !binimg.IsPowerOf2Image(scanner) {
		return nil, errors.New("image must be a square with power-of-2 dimensions")
	}

	if resolution < 1 {
		return nil, errors.New("resolution must be greater than 0")
	}

	// To ensure a consistent behavior and eliminate corner cases,
	// the Quadtree's root node needs to have children. Thus, the
	// first instantiated CNQNode needs to always be subdivided.
	// This condition asserts the resolution is respected.
	if scanner.Bounds().Dx() < resolution*2 {
		return nil, errors.New("the image size must be greater or equal to twice the resolution")
	}

	q := &CNQuadtree{
		resolution: resolution,
		scanner:    scanner,
	}

	q.root = q.newNode(q.scanner.Bounds(), nil, rootQuadrant)
	q.subdivide(q.root)
	return q, nil
}

func (q *CNQuadtree) newNode(bounds image.Rectangle, parent *CNQNode, location quadrant) *CNQNode {
	n := &CNQNode{
		qnode: qnode{
			color:  Gray,
			bounds: bounds,
			parent: parent,
		},
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

func (q *CNQuadtree) subdivide(p *CNQNode) {
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
	nw := q.newNode(image.Rect(x0, y0, x1, y1), p, northWest)
	ne := q.newNode(image.Rect(x1, y0, x2, y1), p, northEast)
	sw := q.newNode(image.Rect(x0, y1, x1, y2), p, southWest)
	se := q.newNode(image.Rect(x1, y1, x2, y2), p, southEast)

	// each sub-quadrant first inherit its parent external neighbours
	// and then updates its internal neighbours.
	nw.cn[west] = p.cn[west]   // inherited
	nw.cn[north] = p.cn[north] // inherited
	nw.cn[east] = ne           // set for decomposition, will need to be updated after
	nw.cn[south] = sw          // set for decomposition, will need to be updated after
	ne.cn[west] = nw           // set for decomposition, will need to be updated after
	ne.cn[north] = p.cn[north] // inherited
	ne.cn[east] = p.cn[east]   // inherited
	ne.cn[south] = se          // set for decomposition, will need to be updated after
	sw.cn[west] = p.cn[west]   // inherited
	sw.cn[north] = nw          // set for decomposition, will need to be updated after
	sw.cn[east] = se           // set for decomposition, will need to be updated after
	sw.cn[south] = p.cn[south] // inherited
	se.cn[west] = sw           // set for decomposition, will need to be updated after
	se.cn[north] = ne          // set for decomposition, will need to be updated after
	se.cn[east] = p.cn[east]   // inherited
	se.cn[south] = p.cn[south] // inherited

	p.northWest = nw
	p.northEast = ne
	p.southWest = sw
	p.southEast = se

	p.updateNECardinalNeighbours()
	p.updateSWCardinalNeighbours()

	// Step3: Updating all neighbours accordingly
	// After the decomposition of a quadrant, all its neighbors in
	// the four directions must be informed of the change so that
	// they can update their own cardinal neighbors accordingly

	// On each direction, a full traversal of the neighbors should
	// be performed. In every quadrant where a reference to the
	// parent quadrant is stored as the Cardinal Neighbor, it
	// should be replaced by one of its children created after the
	// decomposition
	p.updateNeighbours()

	if !nw.isLeaf() {
		q.subdivide(nw)
	}
	if !ne.isLeaf() {
		q.subdivide(ne)
	}
	if !sw.isLeaf() {
		q.subdivide(sw)
	}
	if !se.isLeaf() {
		q.subdivide(se)
	}
}

// Root returns the quadtree root node.
func (q *CNQuadtree) Root() QNode {
	return q.root
}

// ForEachLeaf calls the given function for each leaf node of the quadtree.
//
// Successive calls to the provided function are performed in no particular
// order. The color parameter allows to loop on the leaves of a particular
// color, Black or White.
// NOTE: As by definition, Gray leaves do not exist, passing Gray to
// ForEachLeaf should return all leaves, independently of their color.
func (q *CNQuadtree) ForEachLeaf(color QNodeColor, fn func(QNode)) {
	//panic(" TEST if its workingq.leaves not implemented yet for CNQuadtree")
	for _, n := range q.leaves {
		if color == Gray || n.Color() == color {
			fn(n)
		}
	}
}
