package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type NumberNode struct {
	Pos     errors.Position
	Content *LiteralNode
}

func (this *NumberNode) Name() string {
	return "number"
}

func (this *NumberNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *NumberNode) GetChildren() []Node {
	return []Node{this.Content}
}
