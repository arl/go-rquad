// Code generated by "stringer -type=quadrant"; DO NOT EDIT

package quadtree

import "fmt"

const _quadrant_name = "northWestnorthEastsouthWestsouthEast"

var _quadrant_index = [...]uint8{0, 9, 18, 27, 36}

func (i quadrant) String() string {
	if i < 0 || i >= quadrant(len(_quadrant_index)-1) {
		return fmt.Sprintf("quadrant(%d)", i)
	}
	return _quadrant_name[_quadrant_index[i]:_quadrant_index[i+1]]
}