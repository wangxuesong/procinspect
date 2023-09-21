package semantic

import (
	"bytes"
	"encoding/gob"
)

type (
	NodeEncoder         struct{}
	NodeDecoder[T Node] struct {
		node T
	}
)

func NewNodeEncoder() *NodeEncoder {
	register()
	return &NodeEncoder{}
}

func (n *NodeEncoder) Encode(node Node) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(node)
	return buf.Bytes(), err
}

func NewNodeDecoder[T Node]() *NodeDecoder[T] {
	register()
	return &NodeDecoder[T]{}
}

func (n *NodeDecoder[T]) Decode(data []byte) (Node, error) {
	var node T
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&node)
	return node, err
}
