package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type ImportNode struct {
	Pos    errors.Position
	Source Node
}

func (this *ImportNode) Name() string {
	return "import"
}

func (this *ImportNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *ImportNode) GetChildren() []Node {
	return []Node{this.Source}
}
