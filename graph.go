package quadtree

// TODO: we want to generate an undirected graph where the nodes are the white
// quadnodes and the edges between the nodes indicates that those nodes are
// neighbours

// NodeList is a slice of Node's.
type NodeList []*Node

type Node struct {
	QNode
	links NodeList
}

func newNode(qn QNode) *Node {
	return &Node{
		QNode: qn,
	}
}

type Edge struct {
	n1, n2 *Node
}

type Graph struct {
	nodes NodeList
	edges []Edge
}

func NewGraphFromQuadtree(q Quadtree) *Graph {
	whiteNodes := q.WhiteNodes()
	g := &Graph{
		nodes: make(NodeList, 0, len(whiteNodes)),
	}

	// init a lookup table for fast retrieving of the already created Node's
	nodesLut := map[QNode]*Node{}

	// range over the quadtree nodes
	for _, qn := range whiteNodes {

		// add node to the lookup table
		n := newNode(qn)
		nodesLut[qn] = n

		// save node into the graph
		g.nodes = append(g.nodes, n)

		var nbours QNodeList
		qn.Neighbours(&nbours)
		for _, qnb := range nbours {

			var (
				ok bool
				nb *Node
			)

			// try to take the neighbour from the lookup table
			if nb, ok = nodesLut[qnb]; !ok {
				// neighbour has not been added yet, so create it and add it
				nb = newNode(qnb)
				nodesLut[qnb] = nb
			}

			n.links = append(n.links, nb)
			e := Edge{n, nb}
			g.edges = append(g.edges, e)
		}
	}
	return g
}
