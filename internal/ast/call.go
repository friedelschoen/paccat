package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type CallNode struct {
	Pos    errors.Position
	Target Node
	Args   map[string]Node
}

func (this *CallNode) String() string {
	return fmt.Sprintf("RecipeCall{%v}(%v)", this.Target, this.Args)
}

func (this *CallNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("call"))
	this.Target.WriteHash(hash)
	for key, value := range this.Args {
		hash.Write([]byte(key))
		value.WriteHash(hash)
	}
}

func (this *CallNode) GetPosition() errors.Position {
	return this.Pos
}
