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
	resolution  int            // maximal resolution
	scanner     binimg.Scanner // reference image
	root        *CNQNode       // root node
	whiteNodes  QNodeList      // white nodes (filled during creation)
	onWhiteNode func(QNode)    // callback that fills whiteNodes
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

	// onWhiteNode callback serves the purpose of filling the list of white
	// nodes for providing it to the called of WhiteNodes()
	q.onWhiteNode = func(n QNode) {
		q.whiteNodes = append(q.whiteNodes, n)
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
			// white node callback
			if q.onWhiteNode != nil {
				q.onWhiteNode(n)
			}
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

	// as decomposition is performed in Z-order, the western and
	// northern cardinal neighbours of the parent may not be
	// leaves, if the quadrant in which they lie has been decomposed
	// so we look up for those nodes.
	leafcn0 := p.cn0
	leafcn1 := p.cn1

	// each sub-quadrant first inherit its parent external neighbours
	// and then updates its internal neighbours.
	nw.cn0 = leafcn0 // inherited
	nw.cn1 = leafcn1 // inherited
	nw.cn2 = ne      // set for decomposition, will need to be updated after
	nw.cn3 = sw      // set for decomposition, will need to be updated after
	ne.cn0 = nw      // set for decomposition, will need to be updated after
	ne.cn1 = leafcn1 // inherited
	ne.cn2 = p.cn2   // inherited
	ne.cn3 = se      // set for decomposition, will need to be updated after
	sw.cn0 = leafcn0 // inherited
	sw.cn1 = nw      // set for decomposition, will need to be updated after
	sw.cn2 = se      // set for decomposition, will need to be updated after
	sw.cn3 = p.cn3   // inherited
	se.cn0 = sw      // set for decomposition, will need to be updated after
	se.cn1 = ne      // set for decomposition, will need to be updated after
	se.cn2 = p.cn2   // inherited
	se.cn3 = p.cn3   // inherited

	p.northWest = nw
	p.northEast = ne
	p.southWest = sw
	p.southEast = se

	p.updateNECardinalNeighbours()
	p.updateSWCardinalNeighbours()

	// CHECK INVARIANTS
	// since it is not yet decomposed and thus the parents’
	// eastern cardinal neighbor is itself the Eastern CN of the
	// NE and SE child quadrants
	if p.cn2 != nil {
		if p.cn2.size < p.size {
			panic("should not happen")
		}
		if p.cn2 != ne.cn2 || p.cn2 != se.cn2 {
			panic("should not happen")
		}
	}

	// since it is not yet decomposed and thus the parents’
	// southern cardinal neighbor is itself the Southern CN of
	// the SW and SE child quadrants.
	if p.cn3 != nil {
		if p.cn3.size < p.size {
			panic("should not happen")
		}
		if p.cn3 != sw.cn3 || p.cn3 != se.cn3 {
			panic("should not happen")
		}
	}

	q.updateAllNeighbours(p)

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

// this would be the step 3
func (q *CNQuadtree) updateAllNeighbours(p *CNQNode) {
	// Step3: Updating all neighbours accordingly
	// After the decomposition of a quadrant, all its neighbors in
	// the four directions must be informed of the change so that
	// they can update their own cardinal neighbors accordingly

	// On each direction, a full traversal of the neighbors should
	// be performed. In every quadrant where a reference to the
	// parent quadrant is stored as the Cardinal Neighbor, it
	// should be replaced by one of its children created after the
	// decomposition
	if p.cn0 != nil {
		p.Step3UpdateWest()
	}
	if p.cn1 != nil {
		p.Step3UpdateNorth()
	}
	if p.cn2 != nil {
		p.Step3UpdateEast()
	}
	if p.cn3 != nil {
		p.Step3UpdateSouth()
	}
}

// WhiteNodes returns a slice of all the white nodes of the quadtree.
func (q *CNQuadtree) WhiteNodes() QNodeList {
	return q.whiteNodes
}

// Root returns the quadtree root node.
func (q *CNQuadtree) Root() QNode {
	return q.root
}
