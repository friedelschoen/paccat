package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type GetterNode struct {
	Pos       errors.Position
	Target    Node
	Attribute string
}

func (this *GetterNode) String() string {
	return fmt.Sprintf("RecipeGetter#%s{%v}", this.Attribute, this.Target)
}

func (this *GetterNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("getter"))
	this.Target.WriteHash(hash)
	hash.Write([]byte(this.Attribute))
}

func (this *GetterNode) GetPosition() errors.Position {
	return this.Pos
}
