package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type PanicNode struct {
	Pos     errors.Position
	Message Node
}

func (this *PanicNode) Name() string {
	return "panic"
}

func (this *PanicNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("panic"))
	this.Message.WriteHash(hash)
}

func (this *PanicNode) GetPosition() errors.Position {
	return this.Pos
}

func (this *PanicNode) GetChildren() []Node {
	return []Node{this.Message}
}
