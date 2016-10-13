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
	if n.cn[0] != nil {
		scn0 = fmt.Sprintf("%v-%d", n.cn[0].bounds.Min, n.cn[0].size)
	}
	if n.cn[1] != nil {
		scn1 = fmt.Sprintf("%v-%d", n.cn[1].bounds.Min, n.cn[1].size)
	}
	if n.cn[2] != nil {
		scn2 = fmt.Sprintf("%v-%d", n.cn[2].bounds.Min, n.cn[2].size)
	}
	if n.cn[3] != nil {
		scn3 = fmt.Sprintf("%v-%d", n.cn[3].bounds.Min, n.cn[3].size)
	}
	return fmt.Sprintf("[%v-%d-%s|CN ←%v ↑%v →%v ↓%v]", n.bounds.Min, n.size, n.color, scn0, scn1, scn2, scn3)
}
