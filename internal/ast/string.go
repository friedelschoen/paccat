package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type StringNode struct {
	Pos     errors.Position
	Content []Node
}

func (this *StringNode) Name() string {
	return "string"
}

func (this *StringNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *StringNode) GetChildren() []Node {
	return this.Content
}
