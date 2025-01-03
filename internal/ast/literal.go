package ast

import (
	"fmt"
	"hash"
)

type recipeStringLiteral struct {
	pos     Position
	content string
}

func (this *recipeStringLiteral) String() string {
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", this.content)
}

func (this *recipeStringLiteral) Eval(ctx Context) (Value, error) {
	return &StringValue{
		source:       this,
		Content:      this.content,
		StringSource: []StringSource{},
	}, nil
}

func (this *recipeStringLiteral) WriteHash(hash hash.Hash) {
	hash.Write([]byte("literal"))
	hash.Write([]byte(this.content))
}

func (this *recipeStringLiteral) GetPosition() Position {
	return this.pos
}
