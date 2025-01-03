package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type PanicNode struct {
	Pos     errors.Position
	Message Node
}

func (this *PanicNode) String() string {
	return fmt.Sprintf("RecipePanic")
}

func (this *PanicNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("panic"))
	this.Message.WriteHash(hash)
}

func (this *PanicNode) GetPosition() errors.Position {
	return this.Pos
}
