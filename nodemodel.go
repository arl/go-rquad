package rquad

import (
	"image"
)

type ScanAndSetter interface {
	ScanAndSet(*Node)
}

type NodeModel interface {
	ScanAndSetter
	NewRoot() Node
	NewNode(parent Node, location Quadrant, bounds image.Rectangle) Node
}
