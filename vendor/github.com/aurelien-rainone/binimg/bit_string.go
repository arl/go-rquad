package binimg

import "fmt"

func (c Bit) String() string {
	if c.v == 0 {
		return "Black"
	} else if c.v == 255 {
		return "White"
	}
	return fmt.Sprintf("Bit(%d)", c)
}
