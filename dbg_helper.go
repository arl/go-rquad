package rquad

import "fmt"

func (n *CNQNode) String() string {
	var scn0, scn1, scn2, scn3 string
	if n.cn[West] != nil {
		scn0 = fmt.Sprintf("%v-%d", n.cn[West].bounds.Min, n.cn[West].size)
	}
	if n.cn[North] != nil {
		scn1 = fmt.Sprintf("%v-%d", n.cn[North].bounds.Min, n.cn[North].size)
	}
	if n.cn[East] != nil {
		scn2 = fmt.Sprintf("%v-%d", n.cn[East].bounds.Min, n.cn[East].size)
	}
	if n.cn[South] != nil {
		scn3 = fmt.Sprintf("%v-%d", n.cn[South].bounds.Min, n.cn[South].size)
	}
	return fmt.Sprintf("[%v-%d-%s|CN ←%v ↑%v →%v ↓%v]", n.bounds.Min, n.size, n.color, scn0, scn1, scn2, scn3)
}
