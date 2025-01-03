package ast

import (
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type DictNode struct {
	Pos   errors.Position
	Items map[string]Node
}

func (this *DictNode) String() string {
	return "RecipeDict"
}

func (this *DictNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("list"))
	for _, value := range this.Items {
		value.WriteHash(hash)
	}
}

func (this *DictNode) GetPosition() errors.Position {
	return this.Pos
}
