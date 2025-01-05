package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type DictNode struct {
	Pos   errors.Position
	Items LiteralMap
}

func (this *DictNode) Name() string {
	return "dict"
}

func (this *DictNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	this.Items.WriteHash(hash)
}

func (this *DictNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *DictNode) GetChildren() []Node {
	return []Node{this.Items}
}
