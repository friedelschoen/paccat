package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type LambdaNode struct {
	Pos    errors.Position
	Target Node
	Args   map[string]Node
}

func (this *LambdaNode) String() string {
	return fmt.Sprintf("RecipeLambda#{%v}{%v}", this.Args, this.Target)
}

func (this *LambdaNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("lambda"))
	this.Target.WriteHash(hash)
	for key, value := range this.Args {
		hash.Write([]byte(key))
		if value != nil {
			value.WriteHash(hash)
		}
	}
}

func (this *LambdaNode) GetPosition() errors.Position {
	return this.Pos
}
