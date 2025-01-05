package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type StringNode struct {
	Pos     errors.Position
	Content []Node
}

func (this *StringNode) Name() string {
	return "string"
}

func (this *StringNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("string"))
	for _, content := range this.Content {
		content.WriteHash(hash)
	}
}

func (this *StringNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *StringNode) GetChildren() []Node {
	return this.Content
}
