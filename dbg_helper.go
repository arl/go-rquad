package quadtree

import "fmt"

func (n *Node) String() string {
	return fmt.Sprintf("Node {%v|%d links}",
		n.Bounds(), len(n.links))
}

func (n *CNQNode) String() string {
	var scn0, scn1, scn2, scn3 string
	if n.cn[west] != nil {
		scn0 = fmt.Sprintf("%v-%d", n.cn[west].bounds.Min, n.cn[west].size)
	}
	if n.cn[north] != nil {
		scn1 = fmt.Sprintf("%v-%d", n.cn[north].bounds.Min, n.cn[north].size)
	}
	if n.cn[east] != nil {
		scn2 = fmt.Sprintf("%v-%d", n.cn[east].bounds.Min, n.cn[east].size)
	}
	if n.cn[south] != nil {
		scn3 = fmt.Sprintf("%v-%d", n.cn[south].bounds.Min, n.cn[south].size)
	}
	return fmt.Sprintf("[%v-%d-%s|CN ←%v ↑%v →%v ↓%v]", n.bounds.Min, n.size, n.color, scn0, scn1, scn2, scn3)
}
