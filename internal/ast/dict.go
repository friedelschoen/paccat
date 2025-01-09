package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type DictNode struct {
	Pos   errors.Position
	Items LiteralMap
}

func (this *DictNode) Name() string {
	return "dict"
}

func (this *DictNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *DictNode) GetChildren() []Node {
	return []Node{this.Items}
}
