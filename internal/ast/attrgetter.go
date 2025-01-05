package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type GetterNode struct {
	Pos       errors.Position
	Target    Node
	Attribute *LiteralNode
}

func (this *GetterNode) Name() string {
	return "getter"
}

func (this *GetterNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("getter"))
	this.Target.WriteHash(hash)
	this.Attribute.WriteHash(hash)
}

func (this *GetterNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *GetterNode) GetChildren() []Node {
	return []Node{this.Target, this.Attribute}
}
