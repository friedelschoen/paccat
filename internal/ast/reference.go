package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type ReferenceNode struct {
	Pos      errors.Position
	Variable *LiteralNode
}

func (this *ReferenceNode) Name() string {
	return "reference"
}

func (this *ReferenceNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *ReferenceNode) GetChildren() []Node {
	return []Node{this.Variable}
}
