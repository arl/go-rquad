package rquad

import (
	"image"
)

type ScanAndSetter interface {
	ScanAndSet(*Node)
}

type TreeModel interface {
	ScanAndSetter
	NewRoot() Node
	NewNode(parent Node, location Quadrant, bounds image.Rectangle) Node
}
