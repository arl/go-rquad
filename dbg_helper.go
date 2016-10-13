package quadtree

import "fmt"

func (n *Node) String() string {
	return fmt.Sprintf("Node {%v|%d links}",
		n.Bounds(), len(n.links))
}

func (n *qnode) String() string {
	return fmt.Sprintf("(%v %s)", n.bounds, n.color)
}

func (n *CNQNode) String() string {
	var scn0, scn1, scn2, scn3 string
	if n.cn0 != nil {
		scn0 = fmt.Sprintf("%v-%d", n.cn0.bounds.Min, n.cn0.size)
	}
	if n.cn1 != nil {
		scn1 = fmt.Sprintf("%v-%d", n.cn1.bounds.Min, n.cn1.size)
	}
	if n.cn2 != nil {
		scn2 = fmt.Sprintf("%v-%d", n.cn2.bounds.Min, n.cn2.size)
	}
	if n.cn3 != nil {
		scn3 = fmt.Sprintf("%v-%d", n.cn3.bounds.Min, n.cn3.size)
	}
	return fmt.Sprintf("[%v-%d-%s|CN ←%v ↑%v →%v ↓%v]", n.bounds.Min, n.size, n.color, scn0, scn1, scn2, scn3)
}
