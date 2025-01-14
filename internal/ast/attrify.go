package ast

import "friedelschoen.io/paccat/internal/errors"

type AttrifyNode struct {
	Pos    errors.Position
	Target Node
}

func (this *AttrifyNode) Name() string {
	return "attrify"
}

func (this *AttrifyNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *AttrifyNode) GetChildren() []Node {
	return []Node{this.Target}
}
