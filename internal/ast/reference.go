package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type ReferenceNode struct {
	Pos      errors.Position
	Variable *LiteralNode
}

func (this *ReferenceNode) Name() string {
	return "reference"
}

func (this *ReferenceNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("reference"))
	this.Variable.WriteHash(hash)
}

func (this *ReferenceNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *ReferenceNode) GetChildren() []Node {
	return []Node{this.Variable}
}
