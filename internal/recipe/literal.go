package recipe

import (
	"fmt"
	"hash"
)

type recipeStringLiteral struct {
	value string
}

func (this *recipeStringLiteral) String() string {
	return fmt.Sprintf("RecipeStringLiteral#\"%s\"", string(this.value))
}

func (this *recipeStringLiteral) HasOutput() bool {
	return false
}

func (this *recipeStringLiteral) Eval(ctx *Context) (string, error) {
	return string(this.value), nil
}

func (this *recipeStringLiteral) WriteHash(hash hash.Hash) {
	hash.Write([]byte(this.value))
}
