package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type ImportNode struct {
	Pos    errors.Position
	Source Node
}

func (this *ImportNode) Name() string {
	return "import"
}

func (this *ImportNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("import"))
	this.Source.WriteHash(hash)
}

func (this *ImportNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *ImportNode) GetChildren() []Node {
	return []Node{this.Source}
}
