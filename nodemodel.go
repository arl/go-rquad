package rquad

import "image"

type ScanAndSetter interface {
	ScanAndSet(*Node)
}

type NodeModel interface {
	ScanAndSetter
	NewRoot() Node
	NewNode(Node, Quadrant, image.Rectangle) Node
}
