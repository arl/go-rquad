package quadtree

import (
	"fmt"
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
	whiteNodes := q.WhiteNodes()
	g := &Graph{
		nodes: make([]*Node, len(whiteNodes), len(whiteNodes)),
	}

	// lookup table for fast retrieving of the Node's we
	// created from the QNode coming from the Quadtree
	nodesLut := make(map[QNode]*Node, len(whiteNodes))

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

	var nbours QNodeList

	// range over the quadtree nodes
	for i, qn := range whiteNodes {

		// get node from lut or create a new one
		n := newOrGet(qn)

		// save node into the graph
		g.nodes[i] = n

		nbours = nil
		qn.Neighbours(&nbours)

		// allocate the edges and links slices
		n.links = make([]*Node, len(nbours), len(nbours))
		n.edges = make([]*Edge, len(nbours), len(nbours))

		for j, qnb := range nbours {

			// get neighbour from lut or create a new one
			nb := newOrGet(qnb)
			n.links[j] = nb
			if genEdgeFunc != nil {
				n.edges[j] = genEdgeFunc(qn, qnb)
			}
		}
	}
	return g
}

func newNode(qn QNode) *Node { return &Node{QNode: qn} }

func (n *Node) width() float64 {
	return float64(n.Bounds().Dx())
}

func (n *Node) height() float64 {
	return float64(n.Bounds().Dy())
}

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

	x1 := float64(n.Bounds().Min.X) + n.width()/2
	y1 := float64(n.Bounds().Min.Y) + n.height()/2

	x2 := float64(to.Bounds().Min.X) + to.width()/2
	y2 := float64(to.Bounds().Min.Y) + to.height()/2

	a := math.Abs(x1 - x2)
	b := math.Abs(y1 - y2)

	return a*a + b*b
}

func (n *Node) String() string {
	return fmt.Sprintf("Node {%v|%d links}",
		n.Bounds(), len(n.links))
}
