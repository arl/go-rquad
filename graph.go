package quadtree

import (
	"math"

	astar "github.com/beefsack/go-astar"
)

// A Node is a node, or vertex, of a Graph based off the white QNode's of a
// Quadtree.
type Node struct {
	QNode         // underlying quadtree node
	links []*Node // links to neighbours
	edges []*Edge // edges, edges[i] points to links[i]
}

// An Edge represents an edge between 2 Nodes.
type Edge struct{ n1, n2 *Node }

// Graph is a graph whose in-memory representation is an adjacency list.
//
// Given the use case of this implementation an adjacency list is sufficient
// and required:
// - once the graph has been created, no nodes will ever be added or removed
//   from it
// - the graph main purpose is to be ran A* against so for a given node, the
//   adjacent nodes must be retrieved in 0(1).
// - we don't need to query the graph for adjacency between 2 random nodes.
type Graph struct {
	nodes []*Node
}

// GenEdgeFunc creates an Edge from 2 QNode's.
type GenEdgeFunc func(n1 QNode, n2 QNode) *Edge

// NewGraphFromQuadtree creates a Graph where vertices are the white nodes of
// the given quadtree, and edges exist where 2 white nodes are neighbours in 2D
// space (their underlying rectangle share a segment)
func NewGraphFromQuadtree(q Quadtree, genEdgeFunc GenEdgeFunc) *Graph {
	var whites QNodeList
	q.ForEachLeaf(White, func(n QNode) {
		whites = append(whites, n)
	})

	g := &Graph{
		nodes: make([]*Node, len(whites), len(whites)),
	}

	// lookup table for fast retrieving of the Node's we
	// created from the QNode coming from the Quadtree
	nodesLut := make(map[QNode]*Node, len(whites))

	// get from the lookup table or create the node
	newOrGet := func(qn QNode) *Node {
		var (
			ok bool
			nb *Node
		)
		// try to get the node from the lookup table
		if nb, ok = nodesLut[qn]; !ok {
			// didnt't find it, create it and add it to the lut
			nb = newNode(qn)
			nodesLut[qn] = nb
		}
		return nb
	}

	// range over the quadtree nodes
	for i, qn := range whites {

		// get node from lut or create a new one
		n := newOrGet(qn)

		// save node into the graph
		g.nodes[i] = n

		qn.ForEachNeighbour(func(qnb QNode) {
			// get neighbour from lut or create a new one
			nb := newOrGet(qnb)
			n.links = append(n.links, nb)
			if genEdgeFunc != nil {
				n.edges = append(n.edges, genEdgeFunc(qn, qnb))
			}
		})
	}
	return g
}

func newNode(qn QNode) *Node { return &Node{QNode: qn} }

func (n *Node) center() (x, y float64) {
	sum := n.Bounds().Min.Add(n.Bounds().Max)
	return float64(sum.X) / 2, float64(sum.Y) / 2
}

// PathNeighbors returns the neighbors of the Truck
func (n *Node) PathNeighbors() []astar.Pather {
	nodes := make([]astar.Pather, len(n.links), len(n.links))
	for i := range n.links {
		nodes[i] = n.links[i]
	}
	return nodes
}

// PathNeighborCost returns the cost of the tube leading to Truck.
func (n *Node) PathNeighborCost(to astar.Pather) float64 {
	return math.Sqrt(n.squaredDistance(to.(*Node)))
}

// PathEstimatedCost uses Manhattan distance to estimate
// orthogonal distance between non-adjacent nodes.
func (n *Node) PathEstimatedCost(to astar.Pather) float64 {
	x1, y1 := n.center()
	x2, y2 := to.(*Node).center()
	absX := x1 - x2
	if absX < 0 {
		absX = -absX
	}
	absY := y1 - y2
	if absY < 0 {
		absY = -absY
	}
	return absX + absY
}

// squaredDistance returns the squared straight-line
// distance between this node and another.
func (n *Node) squaredDistance(to *Node) float64 {
	//         a    (x1,y1)
	//    .-----------.
	//    |        .-'
	//  b |     .-'
	//    |  .-'
	//    '-'
	// (x2, y2)

	nb := n.Bounds()
	tob := to.Bounds()

	x1 := float64(nb.Min.X) + float64(nb.Dx()/2)
	y1 := float64(nb.Min.Y) + float64(nb.Dy()/2)

	x2 := float64(tob.Min.X) + float64(tob.Dx()/2)
	y2 := float64(tob.Min.Y) + float64(tob.Dy()/2)

	a := math.Abs(x1 - x2)
	b := math.Abs(y1 - y2)

	return a*a + b*b
}
