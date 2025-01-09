package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type CallNode struct {
	Pos    errors.Position
	Target Node
	Args   LiteralMap
}

func (this *CallNode) Name() string {
	return "call"
}

func (this *CallNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *CallNode) GetChildren() []Node {
	return []Node{this.Target, this.Args}
}
