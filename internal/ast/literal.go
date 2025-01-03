package ast

import (
	"fmt"
	"hash"

	"friedelschoen.io/paccat/internal/errors"
)

type LiteralNode struct {
	Pos     errors.Position
	Content string
}

func (this *LiteralNode) String() string {
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", this.Content)
}

func (this *LiteralNode) WriteHash(hash hash.Hash) {
	hash.Write([]byte("literal"))
	hash.Write([]byte(this.Content))
}

func (this *LiteralNode) GetPosition() errors.Position {
	return this.Pos
}
