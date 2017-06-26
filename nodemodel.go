package rquad

import "image"

type NodeModel interface {
	ScanAndSet(*Node)
	NewRoot() Node
	NewNode(Node, Quadrant, image.Rectangle) Node
}
