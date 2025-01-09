package ast

import (
	"friedelschoen.io/paccat/internal/errors"
)

type LambdaNode struct {
	Pos    errors.Position
	Target Node
	Args   LiteralMap
}

func (this *LambdaNode) Name() string {
	return "lambda"
}

func (this *LambdaNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *LambdaNode) GetChildren() []Node {
	return []Node{this.Target, this.Args}

}
