package ast

import (
	"hash"

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

func (this *LambdaNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("lambda"))
	this.Target.WriteHash(hash)
	this.Args.WriteHash(hash)
}

func (this *LambdaNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *LambdaNode) GetChildren() []Node {
	return []Node{this.Target, this.Args}

}
