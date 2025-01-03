package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type ReferenceNode struct {
	Pos  errors.Position
	Name string
}

func (this *ReferenceNode) String() string {
	return fmt.Sprintf("RecipeReference#%s", this.Name)
}

func (this *ReferenceNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("reference"))
	hash.Write([]byte(this.Name))
}

func (this *ReferenceNode) GetPosition() errors.Position {
	return this.Pos
}
